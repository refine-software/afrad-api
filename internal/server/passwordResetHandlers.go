package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

type passwordResetReq struct {
	Email string `json:"email" binding:"required"`
}

// @Summary      Request Password Reset OTP
// @Description  Generates and sends a password reset OTP to the user's email if the account exists and is verified. Limits OTP requests per day.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body  passwordResetReq  true  "User Email"
// @Success      200  {string}  string  "check your email"
// @Failure      400  {object}  utils.APIError  "Bad request or user not verified or email not found"
// @Failure      403  {object}  utils.APIError  "OTP request limit exceeded"
// @Failure      500  {object}  utils.APIError  "Internal server error"
// @Router       /auth/password-reset [post]
func (s *Server) passwordReset(ctx *gin.Context) {
	var req passwordResetReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}
	userRepo := s.db.User()
	localAuthRepo := s.db.LocalAuth()
	passwordRestRepo := s.db.PasswordReset()
	db := s.db.Pool()

	// check if the  exists
	err = userRepo.CheckEmailExistence(ctx, db, req.Email)
	if errors.Is(err, database.ErrNotFound) {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	// get the user_id by requested email
	userID, err := userRepo.GetIDByEmail(ctx, db, req.Email)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user_id")
		utils.Fail(ctx, apiErr, err)
		return
	}

	// check if user is verified
	Verified, err := localAuthRepo.CheckUserVerification(ctx, db, int32(userID))
	fmt.Println(Verified)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "is_verified")
		utils.Fail(ctx, apiErr, err)
		return
	}

	if !Verified {
		utils.Fail(
			ctx,
			utils.ErrBadRequest,
			errors.New("trying to change password while not verified"),
		)
		return
	}

	// count the password reset otp codes
	countOTPs, err := passwordRestRepo.CountOTPCodesPerDay(ctx, db, int32(userID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "password_reset")
		utils.Fail(ctx, apiErr, err)
		return
	}
	fmt.Println(countOTPs)
	if countOTPs > s.env.MaxOTPRequestsPerDay {
		utils.Fail(ctx, utils.ErrForbidden, errors.New("max otp attempts per day"))
		return
	}

	// generate OTP and store it in the database
	otp := utils.GenerateRandomOTP()

	err = passwordRestRepo.Create(ctx, db, &models.PasswordReset{
		OtpCode:   otp,
		ExpiresAt: utils.GetExpTimeAfterMins(s.env.OTPExpInMin),
		UserID:    int32(userID),
	})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "password_reset")
		utils.Fail(ctx, apiErr, err)
		return
	}

	// send OTP via email
	err = auth.SendOtpEmail(req.Email, otp, s.env)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	utils.Success(ctx, "check your email")
}

type PasswordResetConfirmReq struct {
	NewPassword string `json:"newPassword" binding:"required"`
	OTP         string `json:"otp"         binding:"required"`
	Email       string `json:"email"       binding:"required"`
}

// @Summary      Confirm Password Reset
// @Description  Confirms password reset by verifying the OTP and updating the user's password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body  PasswordResetConfirmReq  true  "New password, OTP, and Email"
// @Success      200  {string}  string  "password changed"
// @Failure      400  {object}  utils.APIError  "Bad request, wrong OTP, or missing fields"
// @Failure      401  {object}  utils.APIError  "OTP expired or unauthorized"
// @Failure      500  {object}  utils.APIError  "Internal server error"
// @Router       /auth/password-reset/confirm [post]
func (s *Server) resetPasswordConfirm(ctx *gin.Context) {
	var req PasswordResetConfirmReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	userRepo := s.db.User()
	localAuthRepo := s.db.LocalAuth()
	passwordRestRepo := s.db.PasswordReset()
	db := s.db.Pool()

	// get the user_id by email
	userID, err := userRepo.GetIDByEmail(ctx, db, req.Email)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	// check if requested otp is the same as what we have in the database
	passwordReset, err := passwordRestRepo.Get(ctx, db, int32(userID))
	if passwordReset.OtpCode != req.OTP {
		utils.Fail(ctx, utils.ErrBadRequest, errors.New("wrong OTP"))
		return
	}

	// check if OTP expired
	if time.Now().After(passwordReset.ExpiresAt) {
		utils.Fail(ctx, utils.ErrUnauthorized, errors.New("OTP is expired"))
		return
	}
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	err = s.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		err = localAuthRepo.Update(ctx, tx, &models.LocalAuth{
			UserID:            int32(userID),
			IsAccountVerified: true,
			PasswordHash:      passwordHash,
		})
		if err != nil {
			return err
		}

		err = passwordRestRepo.Update(ctx, tx, int32(userID))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "password")
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, "password changed")
}
