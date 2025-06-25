package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type ColorRepository interface {
	// Create color, required columns: color
	Create(ctx *gin.Context, db Querier, color string) error

	// Get all colors
	GetAll(ctx *gin.Context, db Querier) (*[]models.Color, error)

	// Update color by id
	Update(ctx *gin.Context, db Querier, id int32, color string) error

	// Check color existenece by id
	CheckExistenece(ctx *gin.Context, db Querier, id int32) error

	// Delete color by id
	Delete(ctx *gin.Context, db Querier, id int32) error
}

type colorRepo struct{}

func NewColorRepository() ColorRepository {
	return &colorRepo{}
}

func (r *colorRepo) Create(ctx *gin.Context, db Querier, color string) error {
	query := `
		INSERT INTO colors(color)
		VALUES ($1)
	`

	_, err := db.Exec(ctx, query, color)
	if err != nil {
		return Parse(
			err,
			"Color",
			"Create",
			Constraints{UniqueViolationCode: "color", NotNullViolationCode: "color"},
		)
	}
	return nil
}

func (r *colorRepo) GetAll(ctx *gin.Context, db Querier) (*[]models.Color, error) {
	query := `
		SELECT id, color
		FROM colors
	`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, Parse(err, "Color", "GetAll", make(Constraints))
	}

	defer rows.Close()
	var colors []models.Color
	for rows.Next() {
		var i models.Color
		err = rows.Scan(&i.ID, &i.Color)
		if err != nil {
			return nil, Parse(err, "Color", "GetAll", make(Constraints))
		}
		colors = append(colors, i)
	}

	err = rows.Err()
	if err != nil {
		return nil, Parse(err, "Color", "GetAll", make(Constraints))
	}

	return &colors, nil
}

func (r *colorRepo) Update(ctx *gin.Context, db Querier, id int32, color string) error {
	query := `
		UPDATE colors
		SET color = $2
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id, color)
	if err != nil {
		return Parse(
			err,
			"Color",
			"Update",
			Constraints{UniqueViolationCode: "color", NotNullViolationCode: "color"},
		)
	}
	return nil
}

func (r *colorRepo) CheckExistenece(ctx *gin.Context, db Querier, id int32) error {
	query := `
		SELECT 1 AS exists
		FROM colors
		WHERE id = $1
	`

	var exists int32
	err := db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return Parse(err, "Color", "CheckExistenece", make(Constraints))
	}
	return nil
}

func (r *colorRepo) Delete(ctx *gin.Context, db Querier, id int32) error {
	query := `
		DELETE FROM colors
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id)
	if err != nil {
		return Parse(err, "Color", "Delete", Constraints{ForeignKeyViolationCode: "id"})
	}
	return nil
}
