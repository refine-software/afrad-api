package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type PhoneVerificationCodeRepository interface {
	// Create phone verification code
	// columns reqired: otp_code, user_id, expires_at
	Create(ctx *gin.Context, db Querier, p *models.PhoneVerification) error

	// Update phone verification code,
	// columns reqired: is_used,
	// by user_id.
	Update(ctx *gin.Context, db Querier, p *models.PhoneVerification) error

	// Get phone verification code,
	// by user_id.
	Get(ctx *gin.Context, db Querier, userID int32) (*models.PhoneVerification, error)

	// count how many otp codes does the user have
	CountUserOtpCodes(ctx *gin.Context, db Querier, userID int) (int, error)
}

type phoneVerificationCodeRepo struct{}

func NewPhoneVerificationCodeRepository() PhoneVerificationCodeRepository {
	return &phoneVerificationCodeRepo{}
}

func (r *phoneVerificationCodeRepo) Create(
	ctx *gin.Context,
	db Querier,
	p *models.PhoneVerification,
) error {
	query := `
		INSERT INTO phone_verification_codes(otp_code, user_id, expires_at)
		VALUES ($1, $2, $3);
	`
	_, err := db.Exec(ctx, query, p.OtpCode, p.UserID, p.ExpiresAt)
	if err != nil {
		return Parse(err)
	}

	return nil
}

func (r *phoneVerificationCodeRepo) Update(
	ctx *gin.Context,
	db Querier,
	p *models.PhoneVerification,
) error {
	query := `
		UPDATE phone_verification_codes
		SET is_used = $2 
		WHERE user_id = $1
	`
	_, err := db.Exec(ctx, query, p.UserID, p.IsUsed)
	if err != nil {
		return Parse(err)
	}
	return nil
}

func (r *phoneVerificationCodeRepo) Get(
	ctx *gin.Context,
	db Querier,
	userID int32,
) (*models.PhoneVerification, error) {
	query := `
		SELECT id, otp_code, is_used, expires_at, created_at
		FROM phone_verification_codes
		WHERE user_id = $1;
	`

	var p models.PhoneVerification
	err := db.QueryRow(ctx, query, userID).
		Scan(&p.ID, &p.OtpCode, &p.IsUsed, &p.ExpiresAt, &p.CreatedAt)
	if err != nil {
		return nil, Parse(err)
	}

	return &p, nil
}

func (r *phoneVerificationCodeRepo) CountUserOtpCodes(
	ctx *gin.Context,
	db Querier,
	userID int,
) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM phone_verification_codes
		WHERE user_id = $1;
	`
	var userOtps int
	err := db.QueryRow(ctx, query, userID).Scan(&userOtps)
	if err != nil {
		return 0, Parse(err)
	}
	return userOtps, nil
}
