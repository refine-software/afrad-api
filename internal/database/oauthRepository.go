package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type OAuthRepository interface {
	Create(ctx *gin.Context, db Querier, authProvider *models.AuthProvider) error
}

type oAuthRepo struct{}

func NewOAuthRepository() OAuthRepository {
	return &oAuthRepo{}
}

func (a *oAuthRepo) Create(
	ctx *gin.Context,
	db Querier,
	authProvider *models.AuthProvider,
) error {
	return nil
}
