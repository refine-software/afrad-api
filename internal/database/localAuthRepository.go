package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type LocalAuthRepository interface {
	// create JWT based authentication,
	// columns required: user_id, phone_number, password_hash.
	Create(ctx *gin.Context, db Querier, l *models.LocalAuth) error

	// update user local auth,
	// required columns: is_phone_verified, password_hash,
	// by user_id.
	Update(ctx *gin.Context, db Querier, l *models.LocalAuth) error
}

type localAuthRepo struct{}

func NewLocalAuthRepository() LocalAuthRepository {
	return &localAuthRepo{}
}

func (r *localAuthRepo) Create(ctx *gin.Context, db Querier, l *models.LocalAuth) error {
	query := `
		INSERT INTO local_auth(user_id, phone_number, password_hash)
		VALUES ($1, $2, $3)
	`
	_, err := db.Exec(ctx, query, l.UserID, l.PhoneNumber, l.PasswordHash)
	if err != nil {
		return Parse(err)
	}
	return nil
}

func (r *localAuthRepo) Update(ctx *gin.Context, db Querier, l *models.LocalAuth) error {
	query := `
		UPDATE local_auth
		SET is_account_verified = $2, password_hash = $3
		WHERE user_id = $1
	`
	_, err := db.Exec(ctx, query, l.UserID, l.IsAccountVerified, l.PasswordHash)
	if err != nil {
		return Parse(err)
	}
	return nil
}
