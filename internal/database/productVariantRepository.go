package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type ProductVariantRepository interface {
	GetAllOfProduct(c *gin.Context, db Querier, productID int32) ([]ProductVariantDetails, error)
	GetPriceByID(c *gin.Context, db Querier, id int32) (int, error)

	// This method will create a product variant.
	//
	// Columns required: quantity, price, product_id, color_id, size_id.
	Create(*gin.Context, Querier, *models.ProductVariant) error
}

type productVariantRepo struct{}

func NewProductVariantRepository() ProductVariantRepository {
	return &productVariantRepo{}
}

type ProductVariantDetails struct {
	ID       int32  `json:"id"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Color    string `json:"color"`
	Size     string `json:"size"`
}

func (pvr *productVariantRepo) GetAllOfProduct(
	c *gin.Context,
	db Querier,
	productID int32,
) ([]ProductVariantDetails, error) {
	query := `
		SELECT 
			pv.id, 
			pv.quantity, 
			pv.price, 
			c.color, 
			s.size || ' (' || s.label || ')' as size
		FROM product_variants pv
		JOIN colors c ON pv.color_id = c.id
		JOIN sizes s ON pv.size_id = s.id
		WHERE pv.product_id = $1
	`

	rows, err := db.Query(c, query, productID)
	if err != nil {
		return nil, Parse(err, "Product Variant", "GetAllOfProduct", make(Constraints))
	}
	defer rows.Close()

	var pvs []ProductVariantDetails
	for rows.Next() {
		var pv ProductVariantDetails
		err = rows.Scan(
			&pv.ID,
			&pv.Quantity,
			&pv.Price,
			&pv.Color,
			&pv.Size,
		)
		if err != nil {
			return nil, Parse(err, "Product Variant", "GetAllOfProduct", make(Constraints))
		}
		pvs = append(pvs, pv)
	}

	if err = rows.Err(); err != nil {
		return nil, Parse(err, "Product Variant", "GetAllOfProduct", make(Constraints))
	}

	return pvs, nil
}

func (pvr *productVariantRepo) Create(
	c *gin.Context,
	db Querier,
	pv *models.ProductVariant,
) error {
	query := `
		INSERT INTO product_variants (quantity, price, product_id, color_id, size_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.Exec(c, query, pv.Quantity, pv.Price, pv.ProductID, pv.ColorID, pv.SizeID)
	if err != nil {
		return Parse(err, "Product Variant", "Create", Constraints{
			UniqueViolationCode:     "variant",
			ForeignKeyViolationCode: "product_id or color_id or size_id",
		})
	}

	return nil
}

func (pvr *productVariantRepo) GetPriceByID(c *gin.Context, db Querier, id int32) (int, error) {
	query := `
		SELECT price
		FROM product_variants
		WHERE id = $1
	`
	var price int

	err := db.QueryRow(c, query, id).Scan(&price)
	if err != nil {
		return 0, Parse(err, "Product Variant", "GetPriceByID", make(Constraints))
	}
	return price, nil
}
