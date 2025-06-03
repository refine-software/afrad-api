package database

import (
	"github.com/gin-gonic/gin"
)

type PasswordResetRepository interface {
	Create(ctx *gin.Context, db Querier) error
}

type passwordResetRepo struct{}

func NewPasswordResetRepository() PasswordResetRepository {
	return &sessionRepo{}
}

func (s *passwordResetRepo) Create(ctx *gin.Context, db Querier) error {
	return nil
}
