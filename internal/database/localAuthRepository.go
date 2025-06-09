package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type LocalAuthRepository interface {
	// fetch the local auth model
	Get(ctx *gin.Context, db Querier, userID int32) (models.LocalAuth, error)

	// create JWT based authentication,
	// columns required: user_id, password_hash.
	Create(ctx *gin.Context, db Querier, l *models.LocalAuth) error

	// update user local auth,
	// required columns: is_phone_verified, password_hash,
	// by user_id.
	Update(ctx *gin.Context, db Querier, l *models.LocalAuth) error

	// update is_account_verified column to true,
	// by user_id.
	UpdateIsAccountVerifiedToTrue(ctx *gin.Context, db Querier, userID int32) error
}

type localAuthRepo struct{}

func NewLocalAuthRepository() LocalAuthRepository {
	return &localAuthRepo{}
}

func (r *localAuthRepo) Get(ctx *gin.Context, db Querier, userID int32) (models.LocalAuth, error) {
	query := `
	SELECT user_id, is_account_verified, password_hash
	FROM local_auth
	WHERE user_id = $1
	`

	var l models.LocalAuth
	err := db.QueryRow(ctx, query, userID).Scan(&l.UserID, &l.IsAccountVerified, &l.PasswordHash)
	if err != nil {
		return models.LocalAuth{}, Parse(err)
	}

	return l, nil
}

func (r *localAuthRepo) Create(ctx *gin.Context, db Querier, l *models.LocalAuth) error {
	query := `
		INSERT INTO local_auth(user_id, password_hash)
		VALUES ($1, $2)
	`
	_, err := db.Exec(ctx, query, l.UserID, l.PasswordHash)
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

func (r *localAuthRepo) UpdateIsAccountVerifiedToTrue(
	ctx *gin.Context,
	db Querier,
	userID int32,
) error {
	query := `
		UPDATE local_auth
		SET is_account_verified = true
		WHERE user_id = $1
	`
	_, err := db.Exec(ctx, query, userID)
	if err != nil {
		return Parse(err)
	}
	return nil
}
