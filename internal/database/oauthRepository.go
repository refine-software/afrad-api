package database

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthProviderRepository interface {
	Create(ctx *gin.Context)
}

type authProviderRepo struct {
	db *pgxpool.Pool
}

func NewAuthProviderRepo(db *pgxpool.Pool) AuthProviderRepository {
	return &authProviderRepo{db}
}

func (a *authProviderRepo) Create(ctx *gin.Context) {
}
