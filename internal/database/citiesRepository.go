package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type CityRepository interface {
	GetAll(ctx *gin.Context, db Querier) (*[]models.City, error)
}

type cityRepo struct{}

func NewCityRepository() CityRepository {
	return &cityRepo{}
}

func (r *cityRepo) GetAll(ctx *gin.Context, db Querier) (*[]models.City, error) {
	query := `
		SELECT id, city
		FROM cities
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, Parse(err, "City", "GetAll", make(Constraints))
	}
	defer rows.Close()
	var cities []models.City
	for rows.Next() {
		var i models.City
		err = rows.Scan(
			&i.ID,
			&i.City,
		)
		if err != nil {
			return nil, Parse(err, "City", "GetAll", make(Constraints))
		}
		cities = append(cities, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, Parse(err, "City", "GetAll", make(Constraints))
	}

	return &cities, nil
}
