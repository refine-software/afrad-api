package database

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/refine-software/afrad-api/internal/models"
)

type CategoryRepository interface {
	// Create new category,
	// required columns: name, parent_id.
	Create(ctx *gin.Context, db Querier, category *models.Category) (int32, error)

	// Get all categories
	GetAll(ctx *gin.Context, db Querier) (*[]models.Category, error)

	// Delete category by id.
	Delete(ctx *gin.Context, db Querier, id int32) error

	// Update category name by id
	Update(ctx *gin.Context, db Querier, id int32, newName string) error
}

type categoryRepo struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepo{}
}

func (r *categoryRepo) Create(
	ctx *gin.Context,
	db Querier,
	category *models.Category,
) (int32, error) {
	query := `
		INSERT INTO categories(name, parent_id)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int32
	err := db.QueryRow(ctx, query, category.Name, category.ParentID).Scan(&id)
	if err != nil {
		return 0, Parse(err, "Category", "Create", Constraints{
			UniqueViolationCode:     "name",
			ForeignKeyViolationCode: "parent",
			NotNullViolationCode:    "name",
		})
	}

	return id, nil
}

func (r *categoryRepo) GetAll(ctx *gin.Context, db Querier) (*[]models.Category, error) {
	query := `
		SELECT id, name, parent_id
		FROM categories
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, Parse(err, "Category", "GetAll", make(Constraints))
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
			return nil, Parse(err, "Category", "GetAll", make(Constraints))
		}
		categories = append(categories, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, Parse(err, "Category", "GetAll", make(Constraints))
	}

	return &categories, nil
}

func (r *categoryRepo) Delete(ctx *gin.Context, db Querier, id int32) error {
	query := `
		DELETE FROM categories 
		WHERE id = $1
	`
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return Parse(err, "Category", "Delete", Constraints{
			ForeignKeyViolationCode: "category", // Could be product or subcategory depending on FK
		})
	}
	if result.RowsAffected() == 0 {
		return Parse(pgx.ErrNoRows, "Category", "Delete", make(Constraints))
	}
	return nil
}

func (r *categoryRepo) Update(
	ctx *gin.Context, db Querier, id int32, newName string,
) error {
	query := `
		UPDATE categories
		SET name = $2
		WHERE id = $1
	`

	result, err := db.Exec(ctx, query, id, newName)
	if err != nil {
		return Parse(err, "Category", "Update", Constraints{
			UniqueViolationCode:  "name",
			NotNullViolationCode: "name",
		})
	}
	if result.RowsAffected() == 0 {
		return Parse(pgx.ErrNoRows, "Category", "Update", make(Constraints))
	}
	return nil
}
