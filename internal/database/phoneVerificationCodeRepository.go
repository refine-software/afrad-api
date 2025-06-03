package database

import (
	"github.com/gin-gonic/gin"
)

type PhoneVerificationCodeRepository interface {
	Create(ctx *gin.Context, db Querier) error
}

type phoneVerificationCodeRepo struct{}

func NewPhoneVerificationCodeRepository() PhoneVerificationCodeRepository {
	return &phoneVerificationCodeRepo{}
}

func (s *phoneVerificationCodeRepo) Create(ctx *gin.Context, db Querier) error {
	return nil
}
