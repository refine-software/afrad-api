package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type ImageRepository interface {
	GetAllOfProduct(
		c *gin.Context,
		db Querier,
		productID int32,
	) ([]models.Image, error)
}

type imageRepo struct{}

func NewImageRepository() ImageRepository {
	return &imageRepo{}
}

func (repo *imageRepo) GetAllOfProduct(
	c *gin.Context,
	db Querier,
	productID int32,
) ([]models.Image, error) {
	query := `
		SELECT id, image, low_res_image
		FROM images
		WHERE product_id = $1
	`

	rows, err := db.Query(c, query, productID)
	if err != nil {
		return nil, Parse(err, "Image", "GetAllOfProduct")
	}
	defer rows.Close()

	var imgs []models.Image
	for rows.Next() {
		var img models.Image
		err = rows.Scan(
			&img.ID,
			&img.Image,
			&img.LowResImage,
		)
		if err != nil {
			return nil, Parse(err, "Image", "GetAllOfProduct")
		}
		imgs = append(imgs, img)
	}

	if err = rows.Err(); err != nil {
		return nil, Parse(err, "Image", "GetAllOfProduct")
	}

	return imgs, nil
}
