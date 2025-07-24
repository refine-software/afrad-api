package database

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/refine-software/afrad-api/internal/models"
)

type SizeRepository interface {
	// Create new size,
	// required columns: size, label.
	Create(ctx *gin.Context, db Querier, size *models.Size) error

	// Get all sizes
	GetAll(ctx *gin.Context, db Querier) (*[]models.Size, error)

	// Get size by label
	GetByLabel(ctx *gin.Context, db Querier, label []string) (*[]models.Size, error)

	// Update size and label by id,
	// required columns:size, label.
	Update(ctx *gin.Context, db Querier, size *models.Size) error

	// Delete size by id
	Delete(ctx *gin.Context, db Querier, id int32) error
}

type sizeRepo struct{}

func NewSizeRepository() SizeRepository {
	return &sizeRepo{}
}

func (r *sizeRepo) Create(ctx *gin.Context, db Querier, size *models.Size) error {
	query := `
		INSERT INTO sizes(size, label)
		VALUES ($1, $2)
	`
	_, err := db.Exec(ctx, query, size.Size, size.Label)
	fmt.Println(size)
	if err != nil {
		return Parse(
			err,
			"Sizes",
			"Create",
			Constraints{UniqueViolationCode: "size, label", NotNullViolationCode: "size, label"},
		)
	}
	return nil
}

func (r *sizeRepo) GetAll(ctx *gin.Context, db Querier) (*[]models.Size, error) {
	query := `
		SELECT id, size, label 
		FROM sizes
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, Parse(err, "Size", "GetAll", make(Constraints))
	}
	defer rows.Close()
	var sizes []models.Size
	for rows.Next() {
		var i models.Size
		err = rows.Scan(
			&i.ID,
			&i.Size,
			&i.Label,
		)
		if err != nil {
			return nil, Parse(err, "Size", "GetAll", make(Constraints))
		}
		sizes = append(sizes, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, Parse(err, "Size", "GetAll", make(Constraints))
	}

	return &sizes, nil
}

// this function assums that there is at least one query or more
func getWhereClauses(queries []string) string {
	var whereClauses []string
	for i := range queries {
		whereClauses = append(whereClauses, fmt.Sprintf("label = $%d", i+1))
	}
	return strings.Join(whereClauses, " OR ")
}

func (r *sizeRepo) GetByLabel(
	ctx *gin.Context,
	db Querier,
	labels []string,
) (*[]models.Size, error) {
	query := fmt.Sprintf(`SELECT id, size, label 
		FROM sizes
		WHERE %s`, getWhereClauses(labels))

	args := make([]any, len(labels))
	for i, v := range labels {
		args[i] = v
	}
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, Parse(err, "Size", "GetByLabel", make(Constraints))
	}

	defer rows.Close()
	var sizes []models.Size
	for rows.Next() {
		var i models.Size
		err = rows.Scan(
			&i.ID,
			&i.Size,
			&i.Label,
		)
		if err != nil {
			return nil, Parse(err, "Size", "GetByLabel", make(Constraints))
		}
		sizes = append(sizes, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, Parse(err, "Size", "GetByLabel", make(Constraints))
	}
	return &sizes, nil
}

func (r *sizeRepo) Update(ctx *gin.Context, db Querier, size *models.Size) error {
	query := `
		UPDATE sizes
		SET
 			size = COALESCE(NULLIF($2, ''), size),
    	label = COALESCE(NULLIF($3, ''), label)
		WHERE id = $1
	`

	result, err := db.Exec(ctx, query, size.ID, size.Size, size.Label)
	if err != nil {
		return Parse(err, "Size", "Update", Constraints{UniqueViolationCode: "size, label"})
	}
	if result.RowsAffected() == 0 {
		return Parse(pgx.ErrNoRows, "Size", "Update", make(Constraints))
	}
	return nil
}

func (r *sizeRepo) Delete(ctx *gin.Context, db Querier, id int32) error {
	query := `
		DELETE FROM sizes
		WHERE id = $1
	`

	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return Parse(err, "Size", "Delete", Constraints{ErrForeignKey: "id"})
	}
	if result.RowsAffected() == 0 {
		return Parse(pgx.ErrNoRows, "Size", "Delete", make(Constraints))
	}

	return nil
}
