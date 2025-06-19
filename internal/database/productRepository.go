package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/refine-software/afrad-api/internal/utils/filters"
)

type ProductRepository interface {
	GetAll(
		ctx *gin.Context,
		db Querier,
		filters filters.Filters,
		prodFilter *filters.ProductFilterOptions,
	) ([]Product, filters.Metadata, error)

	Get(ctx *gin.Context, db Querier, productID int) (*ProductDetails, error)
}

type productRepo struct{}

func NewProductRepository() ProductRepository {
	return &productRepo{}
}

type Product struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	Thumbnail string  `json:"thumbnail"`
	Brand     string  `json:"brand"`
	Category  string  `json:"category"`
	Price     int     `json:"price"`
	Rating    float32 `json:"rating"`
}

func (p *productRepo) GetAll(
	ctx *gin.Context,
	db Querier,
	f filters.Filters,
	productFilters *filters.ProductFilterOptions,
) ([]Product, filters.Metadata, error) {
	whereClause, args := productFilters.GetWhereClause()
	query := fmt.Sprintf(`
	SELECT 
  	COUNT(*) OVER() AS total_records,
  	products.id,
  	products.name,
  	products.thumbnail,
  	brands.brand,
  	categories.name AS category,
  	MIN(product_variants.price) AS min_price,
  	COALESCE(ROUND(AVG(DISTINCT rating_review.rating)::numeric, 2), 0.00) AS avg_rating
	FROM products
	JOIN categories ON categories.id = products.product_category
	JOIN brands ON brands.id = products.brand_id
	JOIN product_variants ON product_variants.product_id = products.id
	LEFT JOIN rating_review ON rating_review.product_id = products.id
	%s
	GROUP BY 
  	products.id,
  	products.name,
  	products.thumbnail,
  	brands.brand,
  	categories.name
	ORDER BY %s %s, products.id ASC
	LIMIT $1 OFFSET $2
	`, whereClause, f.SortColumn(), f.SortDirection())

	fullArgs := []any{f.Limit(), f.Offset()}
	fullArgs = append(fullArgs, args...)
	rows, err := db.Query(ctx, query, fullArgs...)
	if err != nil {
		return nil, filters.Metadata{}, Parse(err, "Product", "GetAll")
	}
	defer rows.Close()

	var (
		totalRecords int
		products     []Product
	)
	for rows.Next() {
		var p Product
		if err = rows.Scan(&totalRecords, &p.ID, &p.Name, &p.Thumbnail, &p.Brand, &p.Category, &p.Price, &p.Rating); err != nil {
			return nil, filters.Metadata{}, Parse(err, "Product", "GetAll")
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, filters.Metadata{}, Parse(err, "Product", "GetAll")
	}

	metadata := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)

	return products, metadata, nil
}

type ProductDetails struct {
	ID         int32       `json:"id"`
	Name       string      `json:"name"`
	Details    pgtype.Text `json:"details"`
	Thumbnail  string      `json:"thumbnail"`
	BrandID    int         `json:"brandId"`
	Brand      string      `json:"brand"`
	CategoryID int         `json:"categoryId"`
	Category   string      `json:"category"`
}

func (pr *productRepo) Get(
	ctx *gin.Context,
	db Querier,
	productID int,
) (*ProductDetails, error) {
	query := `
		SELECT 
			p.id,
			p.name,
			p.details,
			p.thumbnail,
			p.brand_id,
			b.brand,
			p.product_category,
			c.name as category
		FROM products p
		JOIN brands b ON p.brand_id = b.id
		JOIN categories c ON p.product_category = c.id
		WHERE p.id = $1
	`

	var p ProductDetails
	err := db.QueryRow(ctx, query, productID).
		Scan(&p.ID, &p.Name, &p.Details, &p.Thumbnail, &p.BrandID, &p.Brand, &p.CategoryID, &p.Category)
	if err != nil {
		return nil, Parse(err, "Product", "Get")
	}

	return &p, nil
}
