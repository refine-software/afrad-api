package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type PasswordResetRepository interface {
	// Create password reset,
	// required columns: otp_code, expires_at, user_id
	Create(ctx *gin.Context, db Querier, p *models.PasswordReset) *DBError

	// Get password reset by user_id
	Get(ctx *gin.Context, db Querier, userID int32) (*models.PasswordReset, *DBError)

	// update is_used column to true by user_id
	Update(ctx *gin.Context, db Querier, userID int32) *DBError

	// count the OTP codes per day
	CountOTPCodesPerDay(ctx *gin.Context, db Querier, userID int32) (int, *DBError)
}

type passwordResetRepo struct{}

func NewPasswordResetRepository() PasswordResetRepository {
	return &passwordResetRepo{}
}

func (r *passwordResetRepo) Create(ctx *gin.Context, db Querier, p *models.PasswordReset) *DBError {
	query := `
		INSERT INTO password_resets(otp_code, expires_at, user_id)
		VALUES ($1, $2, $3)
	`
	_, err := db.Exec(ctx, query, p.OtpCode, p.ExpiresAt, p.UserID)
	if err != nil {
		return Parse(err, "Password Reset", "Create")
	}
	return nil
}

func (r *passwordResetRepo) Get(
	ctx *gin.Context,
	db Querier,
	userID int32,
) (*models.PasswordReset, *DBError) {
	query := `
		SELECT otp_code, is_used, expires_at
		FROM password_resets
		WHERE user_id = $1
		ORDER BY id DESC
		LIMIT 1
	`
	var p models.PasswordReset
	err := db.QueryRow(ctx, query, userID).Scan(&p.OtpCode, &p.IsUsed, &p.ExpiresAt)
	if err != nil {
		return nil, Parse(err, "Password Reset", "Get")
	}

	return &p, nil
}

func (r *passwordResetRepo) Update(ctx *gin.Context, db Querier, userID int32) *DBError {
	query := `
		UPDATE password_resets
		SET is_used = true
		WHERE user_id = $1
		AND id IN (
			SELECT MAX(id) FROM password_resets
		)
	`
	_, err := db.Exec(ctx, query, userID)
	if err != nil {
		return Parse(err, "Password Reset", "Update")
	}
	return nil
}

func (r *passwordResetRepo) CountOTPCodesPerDay(
	ctx *gin.Context,
	db Querier,
	userID int32,
) (int, *DBError) {
	query := `
	SELECT COUNT(*)
	FROM password_resets
	WHERE user_id = $1
  AND created_at::date = CURRENT_DATE;
	`

	var count int
	err := db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, Parse(err, "Password Reset", "CountOTPCodesPerDay")
	}

	return count, nil
}
