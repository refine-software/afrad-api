package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type UserRepository interface {
	CreateUser(ctx *gin.Context, db Querier, user *models.User) error
	GetUserByEmail(ctx *gin.Context, db Querier, email string) (*models.User, error)
}

type userRepo struct{}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) CreateUser(ctx *gin.Context, db Querier, user *models.User) error {
	return nil
}

func (r *userRepo) GetUserByEmail(
	ctx *gin.Context,
	db Querier,
	email string,
) (*models.User, error) {
	return nil, nil
}
