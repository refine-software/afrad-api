package database

import (
	"github.com/gin-gonic/gin"
)

type SessionRepository interface {
	Create(ctx *gin.Context, db Querier) error
}

type sessionRepo struct{}

func NewSessionRepository() SessionRepository {
	return &sessionRepo{}
}

func (s *sessionRepo) Create(ctx *gin.Context, db Querier) error {
	return nil
}
