package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type UserRepository interface {
	// This method will create the user, with the following data:
	// first_name, last_name, image, email, phone_number, role.
	Create(ctx *gin.Context, db Querier, user *models.User) (int, *DBError)

	// This method will update the following user columns:
	// first_name, last_name, image, role.
	// based on the user id.
	Update(ctx *gin.Context, db Querier, u *models.User) *DBError

	// Get user by id
	Get(ctx *gin.Context, db Querier, id int) (*models.User, *DBError)

	// Get user by email
	GetByEmail(ctx *gin.Context, db Querier, email string) (*models.User, *DBError)

	// Get the user role by user id
	GetRole(ctx *gin.Context, db Querier, id int32) (models.Role, *DBError)

	// Get user id by email
	GetIDByEmail(ctx *gin.Context, db Querier, email string) (int, *DBError)

	// Check if email is exist
	CheckEmailExistence(ctx *gin.Context, db Querier, email string) *DBError
}

type userRepo struct{}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) Create(ctx *gin.Context, db Querier, u *models.User) (int, *DBError) {
	query := `
	INSERT INTO users(first_name, last_name, image, email, phone_number, role)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	var id int
	err := db.QueryRow(ctx, query, u.FirstName, u.LastName, u.Image, u.Email, u.PhoneNumber, u.Role).
		Scan(&id)
	if err != nil {
		return 0, Parse(err, "User", "Create")
	}

	return id, nil
}

func (r *userRepo) Update(ctx *gin.Context, db Querier, u *models.User) *DBError {
	query := `
	UPDATE users
	SET first_name = $1, last_name = $2, image = $3, role = $4
	WHERE id = $5;`

	_, err := db.Exec(ctx, query, u.FirstName, u.LastName, u.Image, u.Role, u.ID)
	if err != nil {
		return Parse(err, "User", "Update")
	}

	return nil
}

func (r *userRepo) Get(
	ctx *gin.Context,
	db Querier,
	id int,
) (*models.User, *DBError) {
	query := `SELECT id, first_name, last_name, image, email, role, created_at, updated_at, phone_number
	FROM users 
	WHERE id = $1
	`
	var u models.User
	err := db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Image,
		&u.Email,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.PhoneNumber,
	)
	if err != nil {
		return nil, Parse(err, "User", "Get")
	}

	return &u, nil
}

func (r *userRepo) GetByEmail(
	ctx *gin.Context,
	db Querier,
	email string,
) (*models.User, *DBError) {
	query := `SELECT id, first_name, last_name, image, email, role, created_at, updated_at
	FROM users 
	WHERE email = $1
	`
	var u models.User

	err := db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Image,
		&u.Email,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, Parse(err, "User", "GetByEmail")
	}

	return &u, nil
}

func (r *userRepo) GetRole(ctx *gin.Context, db Querier, id int32) (models.Role, *DBError) {
	query := `SELECT role
	FROM users
	WHERE id = $1
	`

	var role models.Role
	err := db.QueryRow(ctx, query, id).Scan(&role)
	if err != nil {
		return "", Parse(err, "User", "GetRole")
	}

	return role, nil
}

func (r *userRepo) GetIDByEmail(ctx *gin.Context, db Querier, email string) (int, *DBError) {
	query := `
	SELECT id
	FROM users
	WHERE email = $1
	`
	var id int

	err := db.QueryRow(ctx, query, email).Scan(&id)
	if err != nil {
		return 0, Parse(err, "User", "GetIDByEmail")
	}

	return id, nil
}

func (r *userRepo) CheckEmailExistence(ctx *gin.Context, db Querier, email string) *DBError {
	query := `
		SELECT 1 AS exist FROM users
		WHERE email = $1
	`
	var exist int32
	err := db.QueryRow(ctx, query, email).Scan(&exist)
	if err != nil {
		return Parse(err, "User", "CheckEmailExistence")
	}
	return nil
}
