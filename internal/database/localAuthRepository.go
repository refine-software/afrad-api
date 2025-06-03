package database

import "github.com/gin-gonic/gin"

type LocalAuthRepository interface {
	CreateLocalAuth(ctx *gin.Context, db Querier) error
}

type localAuthRepo struct{}

func NewLocalAuthRepository() LocalAuthRepository {
	return &localAuthRepo{}
}

func (l *localAuthRepo) CreateLocalAuth(ctx *gin.Context, db Querier) error {
	return nil
}
