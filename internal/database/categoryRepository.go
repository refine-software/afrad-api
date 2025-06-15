package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type CategoryRepository interface {
	Create(ctx *gin.Context, db Querier, category *models.Category) (int32, *DBError)
	GetAll(ctx *gin.Context, db Querier) (*[]models.Category, *DBError)
}

type categoryRepo struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepo{}
}

func (r *categoryRepo) Create(
	ctx *gin.Context,
	db Querier,
	category *models.Category,
) (int32, *DBError) {
	query := `
		INSERT INTO categories(name, parent_id)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int32
	err := db.QueryRow(ctx, query, category.Name, category.ParentID).Scan(&id)
	if err != nil {
		return 0, Parse(err, "Category", "Create")
	}

	return id, nil
}

func (r *categoryRepo) GetAll(ctx *gin.Context, db Querier) (*[]models.Category, *DBError) {
	query := `
		SELECT id, name, parent_id
		FROM categories
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, Parse(err, "Category", "Get")
	}
	defer rows.Close()
	var categories []models.Category
	for rows.Next() {
		var i models.Category
		err = rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
		)
		if err != nil {
			return nil, Parse(err, "categoryRepo", "Get")
		}
		categories = append(categories, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, Parse(err, "categoryRepo", "Get")
	}

	return &categories, nil
}
