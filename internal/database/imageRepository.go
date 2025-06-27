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

	// This method will create a record in the images table.
	//
	// Columns required: image, low_res_image, product_id.
	Create(*gin.Context, Querier, *models.Image) error
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
		return nil, Parse(err, "Image", "GetAllOfProduct", make(Constraints))
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
			return nil, Parse(err, "Image", "GetAllOfProduct", make(Constraints))
		}
		imgs = append(imgs, img)
	}

	if err = rows.Err(); err != nil {
		return nil, Parse(err, "Image", "GetAllOfProduct", make(Constraints))
	}

	return imgs, nil
}

func (repo *imageRepo) Create(
	c *gin.Context,
	db Querier,
	i *models.Image,
) error {
	query := `
		INSERT INTO images (image, low_res_image, product_id)
		VALUES ($1, $2, $3)
	`

	_, err := db.Exec(c, query, i.Image, i.LowResImage, i.ProductID)
	if err != nil {
		return Parse(err, "Image", "Create", Constraints{
			UniqueViolationCode:     "image or low_res_image",
			ForeignKeyViolationCode: "product",
			NotNullViolationCode:    "image or low_res_image or product_id",
		})
	}

	return nil
}
