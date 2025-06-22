package database

import (
	"github.com/gin-gonic/gin"
)

type ProductVariantRepository interface {
	GetAllOfProduct(c *gin.Context, db Querier, productID int32) ([]ProductVariantDetails, error)
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
