package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type UserRepository interface {
	CreateUser(ctx *gin.Context, db Querier, user *models.User) error
	GetUserByEmail(ctx *gin.Context, db Querier, email string) (*models.User, error)
	UpdateUser(ctx *gin.Context, db Querier, u *models.User) error
}

type userRepo struct{}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

// this method will create the user, with the following data:
// first_name, last_name, image, email, role.
func (r *userRepo) CreateUser(ctx *gin.Context, db Querier, u *models.User) error {
	query := `
	INSERT INTO users(first_name, last_name, image, email, role)
	VALUES($1, $2, $3, $4, $5)
	`

	_, err := db.Exec(ctx, query, u.FirstName, u.LastName, u.Image, u.Email, u.Role)
	if err != nil {
		return Parse(err)
	}

	return nil
}

// this method will update the following user columns:
// first_name, last_name, image, role.
// based on the user id.
func (r *userRepo) UpdateUser(ctx *gin.Context, db Querier, u *models.User) error {
	query := `
	UPDATE users
	SET first_name = $1, last_name = $2, image = $3, role = $4
	WHERE id = $5;`

	_, err := db.Exec(ctx, query, u.FirstName, u.LastName, u.Image, u.Role, u.UpdatedAt, u.ID)
	if err != nil {
		return Parse(err)
	}

	return nil
}

func (r *userRepo) GetUserByEmail(
	ctx *gin.Context,
	db Querier,
	email string,
) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	u := &models.User{}

	err := db.QueryRow(ctx, query, email).Scan(u)
	if err != nil {
		return nil, Parse(err)
	}

	return u, nil
}
