package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type AccountVerificationCodeRepository interface {
	// Create phone verification code,
	// columns reqired: otp_code, user_id, expires_at
	Create(ctx *gin.Context, db Querier, a *models.AccountVerificationCode) error

	// Update phone verification code,
	// columns reqired: is_used,
	// by user_id.
	Update(ctx *gin.Context, db Querier, a *models.AccountVerificationCode) error

	// Get phone verification code,
	// by user_id.
	Get(ctx *gin.Context, db Querier, userID int32) (*models.AccountVerificationCode, error)

	// count how many otp codes does the user have
	CountUserOtpCodes(ctx *gin.Context, db Querier, userID int) (int, error)

	// count the number of otp codes a user have in a day
	CountUserOTPCodesPerDay(
		ctx *gin.Context,
		db Querier,
		userID int32,
	) (int, error)
}

type accountVerificationCodeRepo struct{}

func NewAccountVerificationCodeRepository() AccountVerificationCodeRepository {
	return &accountVerificationCodeRepo{}
}

func (r *accountVerificationCodeRepo) Create(
	ctx *gin.Context,
	db Querier,
	p *models.AccountVerificationCode,
) error {
	query := `
		INSERT INTO account_verification_codes(otp_code, user_id, expires_at)
		VALUES ($1, $2, $3);
	`
	_, err := db.Exec(ctx, query, p.OtpCode, p.UserID, p.ExpiresAt)
	if err != nil {
		return Parse(err)
	}

	return nil
}

func (r *accountVerificationCodeRepo) Update(
	ctx *gin.Context,
	db Querier,
	p *models.AccountVerificationCode,
) error {
	query := `
		UPDATE account_verification_codes
		SET is_used = $2 
		WHERE user_id = $1
	`
	_, err := db.Exec(ctx, query, p.UserID, p.IsUsed)
	if err != nil {
		return Parse(err)
	}
	return nil
}

func (r *accountVerificationCodeRepo) Get(
	ctx *gin.Context,
	db Querier,
	userID int32,
) (*models.AccountVerificationCode, error) {
	query := `
		SELECT id, otp_code, is_used, expires_at, created_at
		FROM account_verification_codes
		WHERE user_id = $1;
	`

	var a models.AccountVerificationCode
	err := db.QueryRow(ctx, query, userID).
		Scan(&a.ID, &a.OtpCode, &a.IsUsed, &a.ExpiresAt, &a.CreatedAt)
	if err != nil {
		return nil, Parse(err)
	}

	return &a, nil
}

func (r *accountVerificationCodeRepo) CountUserOtpCodes(
	ctx *gin.Context,
	db Querier,
	userID int,
) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM account_verification_codes
		WHERE user_id = $1;
	`
	var userOtps int
	err := db.QueryRow(ctx, query, userID).Scan(&userOtps)
	if err != nil {
		return 0, Parse(err)
	}
	return userOtps, nil
}

func (r *accountVerificationCodeRepo) CountUserOTPCodesPerDay(
	ctx *gin.Context,
	db Querier,
	userID int32,
) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM account_verification_codes
	WHERE user_id = $1
  AND created_at::date = CURRENT_DATE;
	`

	var count int
	err := db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, Parse(err)
	}

	return count, nil
}
