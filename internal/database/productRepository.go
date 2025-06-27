package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils/filters"
)

type ProductRepository interface {
	GetAll(
		ctx *gin.Context,
		db Querier,
		filters filters.Filters,
		prodFilter *filters.ProductFilterOptions,
	) ([]Product, filters.Metadata, error)

	GetDetails(ctx *gin.Context, db Querier, productID int) (*ProductDetails, error)

	Get(ctx *gin.Context, db Querier, productID int) (*models.Product, error)

	// This function will create a product.
	//
	// Columns required: name, details, thumbnail, brand_id, product_category.
	// Returns: id.
	Create(*gin.Context, Querier, *models.Product) (productID int32, err error)

	// This Method will update the product.
	//
	// Columns required: name, details, brand_id, product_category.
	// By: id.
	Update(*gin.Context, Querier, *models.Product) error
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

func (pr *productRepo) GetAll(
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
		return nil, filters.Metadata{}, Parse(err, "Product", "GetAll", make(Constraints))
	}
	defer rows.Close()

	var (
		totalRecords int
		products     []Product
	)
	for rows.Next() {
		var p Product
		if err = rows.Scan(&totalRecords, &p.ID, &p.Name, &p.Thumbnail, &p.Brand, &p.Category, &p.Price, &p.Rating); err != nil {
			return nil, filters.Metadata{}, Parse(err, "Product", "GetAll", make(Constraints))
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, filters.Metadata{}, Parse(err, "Product", "GetAll", make(Constraints))
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

func (pr *productRepo) GetDetails(
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
		return nil, Parse(err, "Product", "Get", make(Constraints))
	}

	return &p, nil
}

func (pr *productRepo) Get(
	ctx *gin.Context,
	db Querier,
	productID int,
) (*models.Product, error) {
	query := `
		SELECT 
			id,
			name,
			details,
			thumbnail,
			created_at,
			updated_at,
			brand_id,
			product_category,
		FROM products
		WHERE id = $1
	`

	var p models.Product
	err := db.QueryRow(ctx, query, productID).
		Scan(&p.ID, &p.Name, &p.Details, &p.Thumbnail, &p.CreatedAt, &p.UpdatedAt, &p.BrandID, &p.CreatedAt)
	if err != nil {
		return nil, Parse(err, "Product", "Get", make(Constraints))
	}

	return &p, nil
}

func (pr *productRepo) Create(
	c *gin.Context,
	db Querier,
	p *models.Product,
) (productID int32, err error) {
	query := `
		INSERT INTO products (name, details, thumbnail, brand_id, product_category)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = db.QueryRow(c, query, p.Name, p.Details, p.Thumbnail, p.BrandID, p.ProductCategory).
		Scan(&productID)
	if err != nil {
		return 0, Parse(err, "Product", "Create", Constraints{
			UniqueViolationCode:           "name",
			ForeignKeyViolationCode:       "brand_id or product_category", // can’t distinguish which without inspecting error detail
			NotNullViolationCode:          "name or thumbnail or brand_id or product_category",
			StringDataRightTruncationCode: "name or thumbnail",
		})
	}

	return productID, nil
}

func (pr *productRepo) Update(
	c *gin.Context,
	db Querier,
	p *models.Product,
) error {
	query := `
		UPDATE products
		SET name = $2, details = $3, brand_id = $4, product_category = $5
		WHERE id = $1
	`

	_, err := db.Exec(c, query, p.ID, p.Name, p.Details, p.BrandID, p.ProductCategory)
	if err != nil {
		return Parse(err, "Product", "Update", Constraints{
			UniqueViolationCode:     "name",
			ForeignKeyViolationCode: "brand_id or product_category",
			NotNullViolationCode:    "name or brand_id or product_category",
		})
	}

	return nil
}
