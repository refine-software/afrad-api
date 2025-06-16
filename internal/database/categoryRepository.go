package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type CategoryRepository interface {
	// Create new category,
	// required columns: name, parent_id.
	Create(ctx *gin.Context, db Querier, category *models.Category) (int32, *DBError)

	// Get all categories
	GetAll(ctx *gin.Context, db Querier) (*[]models.Category, *DBError)

	// Delete category by id.
	Delete(ctx *gin.Context, db Querier, id int32) *DBError

	// Check if category exists by id.
	CheckExistence(ctx *gin.Context, db Querier, id int32) *DBError

	// Update category name by id
	Update(ctx *gin.Context, db Querier, id int32, newName string) *DBError
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
		return nil, Parse(err, "Category", "GetAll")
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
			return nil, Parse(err, "Category", "GetAll")
		}
		categories = append(categories, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, Parse(err, "Category", "GetAll")
	}

	return &categories, nil
}

func (r *categoryRepo) Delete(ctx *gin.Context, db Querier, id int32) *DBError {
	query := `
		DELETE FROM categories 
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id)
	if err != nil {
		return Parse(err, "Category", "Delete")
	}
	return nil
}

func (r *categoryRepo) CheckExistence(ctx *gin.Context, db Querier, id int32) *DBError {
	query := `
		SELECT 1 AS exist FROM categories
		WHERE id = $1
	`
	var exist int32
	err := db.QueryRow(ctx, query, id).Scan(&exist)
	if err != nil {
		return Parse(err, "Category", "CheckExistence")
	}
	return nil
}

func (r *categoryRepo) Update(
	ctx *gin.Context, db Querier, id int32, newName string,
) *DBError {
	query := `
		UPDATE categories
		SET name = $2
		WHERE id = $1
	`

	_, err := db.Exec(ctx, query, id, newName)
	if err != nil {
		return Parse(err, "Category", "Update")
	}
	return nil
}
