package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

type registerReq struct {
	FirstName   string `form:"firstName"   binding:"required"`
	LastName    string `form:"lastName"    binding:"required"`
	Email       string `form:"email"       binding:"required"`
	PhoneNumber string `form:"phoneNumber"`
	Password    string `form:"password"    binding:"required"`
}

func (s *Server) register(ctx *gin.Context) {
	// get user info
	var req registerReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	// upload image if exists
	var imageURL string
	file, fileHeader, err := ctx.Request.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			file = nil
			fileHeader = nil

		} else {
			utils.Fail(ctx, utils.ErrBadRequest, err)
			return
		}
	} else {
		imageURL, err = s.s3.UploadImage(ctx, file, fileHeader)
		if err != nil {
			utils.Fail(ctx, utils.ErrInternal, err)
			return
		}

	}

	userRepo := s.db.User()
	localAuthRepo := s.db.LocalAuth()
	otpCodeRepo := s.db.PhoneVerificationCode()

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Image:     imageURL,
		Role:      "user",
	}

	// hash tha password
	passwordHashed, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	// start a transaction for creating a user
	err = s.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		var userID int
		userID, err = userRepo.Create(ctx, tx, user)
		if err != nil {
			return err
		}

		err = localAuthRepo.Create(ctx, tx, &models.LocalAuth{
			UserID:       int32(userID),
			PhoneNumber:  req.PhoneNumber,
			PasswordHash: passwordHashed,
		})
		if err != nil {
			return err
		}

		// generate OTP
		otp := utils.GenerateRandomOTP()
		err = otpCodeRepo.Create(ctx, tx, &models.PhoneVerification{
			UserID:    int32(userID),
			OtpCode:   otp,
			ExpiresAt: time.Now().Add(time.Minute * 5),
		})
		if err != nil {
			return err
		}
		err = auth.SendVerificationEmail(req.Email, otp, s.env)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Created(ctx, "user Created", "")
}

type verifyAccountReq struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp"   binding:"required"`
}

func (s *Server) verifyAccount(ctx *gin.Context) {
	var req verifyAccountReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	db := s.db.Pool()
	userRepo := s.db.User()
	otpCodeRepo := s.db.PhoneVerificationCode()
	localAuthRepo := s.db.LocalAuth()

	// check if email exists in the database
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

	// limit the otps to 10
	var otpCodes int
	otpCodes, err = otpCodeRepo.CountUserOtpCodes(ctx, db, userID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "otp")
		utils.Fail(ctx, apiErr, err)
		return
	}

	if otpCodes > 10 {
		utils.Fail(ctx, utils.ErrBadRequest, errors.New("you reach maximam OTP code request"))
		return
	}

	// get otp_code and otp Expires_at by user_id
	otp, err := otpCodeRepo.Get(ctx, db, int32(userID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "otp_code")
		utils.Fail(ctx, apiErr, err)
		return
	}

	// check if requested otp is the same as what we have in the database
	if otp.OtpCode != req.OTP {
		utils.Fail(ctx, utils.ErrBadRequest, errors.New("wrong OTP code, Try again"))
		return
	}

	// check if otp code is expired
	if time.Now().After(otp.ExpiresAt) {
		utils.Fail(ctx, utils.ErrUnauthorized, errors.New("your OTP is expired"))
		return
	}

	// start a transaction
	err = s.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		err = localAuthRepo.Update(ctx, tx, &models.LocalAuth{
			UserID:            int32(userID),
			IsAccountVerified: true,
		})
		if err != nil {
			return err
		}
		err = otpCodeRepo.Update(ctx, tx, &models.PhoneVerification{
			UserID: int32(userID),
			IsUsed: true,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user_verification")
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, "you are verified", nil)
}
