package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

// @Summary      Register User
// @Description  Registers a new user with optional profile image and sends a verification OTP via email.
// @Tags         Auth
// @Accept       multipart/form-data
// @Produce      json
// @Param        firstName    formData  string true  "First Name"
// @Param        lastName     formData  string true  "Last Name"
// @Param        email        formData  string true  "Email"
// @Param        phoneNumber  formData  string false "Phone Number"
// @Param        password     formData  string true  "Password"
// @Param        image        formData  file   false "Optional Profile Image"
// @Success      201  {string}  string  "user created"
// @Failure      400  {object}  utils.APIError  "Invalid request data"
// @Failure      500  {object}  utils.APIError  "Internal server error"
// @Router       /auth/register [post]
func (s *Server) register(ctx *gin.Context) {
	// get user info
	var req registerReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	imgURL := pgtype.Text{String: "", Valid: false}
	imgUpload, apiErr := getImageFile(ctx)
	if apiErr != nil {
		utils.Fail(ctx, apiErr, nil)
		return
	}
	if imgUpload != nil {
		defer imgUpload.File.Close()

		var uploadedURL string
		uploadedURL, err = s.s3.UploadImage(ctx, imgUpload.File, imgUpload.Header)
		if err != nil {
			utils.Fail(ctx, utils.ErrInternal, err)
			return
		}
		imgURL.String = uploadedURL
		imgURL.Valid = true
	}

	userRepo := s.db.User()
	localAuthRepo := s.db.LocalAuth()
	otpCodeRepo := s.db.AccountVerificationCode()

	user := &models.User{
		FirstName:   req.FirstName,
		LastName:    pgtype.Text{String: req.LastName, Valid: true},
		Email:       req.Email,
		PhoneNumber: pgtype.Text{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		Image:       imgURL,
		Role:        getUserRole(req.Email),
	}

	// hash the password
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
			PasswordHash: passwordHashed,
		})
		if err != nil {
			return err
		}

		// generate OTP
		otp := utils.GenerateRandomOTP()
		err = otpCodeRepo.Create(ctx, tx, &models.AccountVerificationCode{
			UserID:    int32(userID),
			OtpCode:   otp,
			ExpiresAt: utils.GetExpTimeAfterMins(s.env.OTPExpInMin),
		})
		if err != nil {
			return err
		}

		err = auth.SendOtpEmail(req.Email, otp, s.env)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(ctx, apiErr, err)
		if imgURL.String != "" {
			_ = s.s3.DeleteImageByURL(ctx, imgURL.String)
		}
		return
	}

	utils.Created(ctx, "user created")
}

type verifyAccountReq struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp"   binding:"required"`
}

// @Summary      Verify Account
// @Description  Verifies a user's account using email and OTP.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body  verifyAccountReq  true  "Verification Data"
// @Success      200  {string}  string  "your account has been verified"
// @Failure      400  {object}  utils.APIError  "Bad request or invalid OTP"
// @Failure      401  {object}  utils.APIError  "OTP expired"
// @Failure      500  {object}  utils.APIError  "Internal server error"
// @Router       /auth/verify-account [post]
func (s *Server) verifyAccount(ctx *gin.Context) {
	var req verifyAccountReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	db := s.db.Pool()
	userRepo := s.db.User()
	otpCodeRepo := s.db.AccountVerificationCode()
	localAuthRepo := s.db.LocalAuth()

	// check if email exists in the database
	err = userRepo.CheckEmailExistence(ctx, db, req.Email)
	if err != nil && database.IsDBNotFoundErr(err) {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "email not found",
			},
			err,
		)
		return
	}

	// get the user_id by requested email
	userID, err := userRepo.GetIDByEmail(ctx, db, req.Email)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user_id")
		utils.Fail(ctx, apiErr, err)
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
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "wrong OTP, try again",
			},
			nil,
		)
		return
	}

	// check if otp code is expired
	if time.Now().After(otp.ExpiresAt) {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusUnauthorized,
				Message: "OTP is expired, try resending a new one",
			},
			nil,
		)
		return
	}

	// start a transaction
	err = s.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		err = localAuthRepo.UpdateIsAccountVerifiedToTrue(ctx, tx, int32(userID))
		if err != nil {
			return err
		}
		err = otpCodeRepo.Update(ctx, tx, &models.AccountVerificationCode{
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

	utils.Success(ctx, "your account has been verified")
}

type resendVerificationOTPReq struct {
	Email string `json:"email" binding:"required"`
}

// @Summary      Resend Verification OTP
// @Description  Resends an OTP code to a user's email if the account is not yet verified.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body  resendVerificationOTPReq  true  "Email for which to resend OTP"
// @Success      200  {string}  string  "check your email for otp"
// @Failure      400  {object}  utils.APIError  "Bad request, invalid input, or already verified"
// @Failure      403  {object}  utils.APIError  "OTP request limit reached"
// @Failure      500  {object}  utils.APIError  "Internal server error"
// @Router       /auth/resend-verification [post]
func (s *Server) resendVerification(c *gin.Context) {
	// get user email
	var req resendVerificationOTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, utils.ErrBadRequest, err)
		return
	}

	accVerificationRepo := s.db.AccountVerificationCode()
	userRepo := s.db.User()
	localAuthRepo := s.db.LocalAuth()
	db, err := s.db.BeginTx(c)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}
	committed := false
	defer func() {
		if p := recover(); p != nil {
			_ = db.Rollback(c)
			panic(p)
		} else if !committed {
			_ = db.Rollback(c)
		}
	}()

	// get user
	user, err := userRepo.GetByEmail(c, db, req.Email)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	localAuth, err := localAuthRepo.Get(c, db, user.ID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	if localAuth.IsAccountVerified {
		utils.Fail(
			c,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "Account is already verified. No need to resend verification email.",
			},
			nil,
		)
		return
	}

	// limit verification otps to 10 per day
	numOfOTPs, err := accVerificationRepo.CountUserOTPCodesPerDay(c, db, user.ID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "otp")
		utils.Fail(c, apiErr, err)
		return
	}

	if numOfOTPs > s.env.MaxOTPRequestsPerDay {
		utils.Fail(
			c,
			&utils.APIError{
				Code:    http.StatusForbidden,
				Message: "You've reached the limit of otp requests",
			},
			nil,
		)
		return
	}

	// generate OTP
	otp := utils.GenerateRandomOTP()

	// store OTP
	a := models.AccountVerificationCode{
		OtpCode:   otp,
		ExpiresAt: utils.GetExpTimeAfterMins(s.env.OTPExpInMin),
		UserID:    user.ID,
	}
	err = accVerificationRepo.Create(c, db, &a)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "otp")
		utils.Fail(c, apiErr, err)
		return
	}

	// send verificaion OTP
	err = auth.SendOtpEmail(user.Email, otp, s.env)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	err = db.Commit(c)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}
	committed = true

	// responed
	utils.Success(c, "check your email for otp")
}

type refreshTokenReq struct {
	UserID int32 `json:"userId"`
}

type refreshTokenRes struct {
	AccessToken string `json:"accessToken"`
}

// @Summary      Refresh Tokens
// @Description  Rotates a valid refresh token and returns a new access token. Requires refresh token in cookie.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body  refreshTokenReq  true  "User ID"
// @Success      200  {object}  refreshTokenRes  "New access token"
// @Failure      400  {object}  utils.APIError  "Bad request or missing refresh_token cookie"
// @Failure      401  {object}  utils.APIError  "Invalid, expired, or revoked session"
// @Failure      500  {object}  utils.APIError  "Internal server error"
// @Router       /auth/refresh [post]
// @Security     RefreshTokenCookie
func (s *Server) refreshTokens(c *gin.Context) {
	var req refreshTokenReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(c, utils.ErrBadRequest, err)
		return
	}

	// Accept a refresh token and a user agent
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.Fail(
			c,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "refresh_token cookie is required",
			},
			err,
		)
		return
	}

	userAgent := getHeader(c, "User-Agent")
	if userAgent == "" {
		return
	}

	sessionRepo := s.db.Session()
	userRepo := s.db.User()
	db := s.db.Pool()

	// Validate the refresh token:
	// 		- exists in the database
	session, err := sessionRepo.GetByUserIDAndUserAgent(c, db, req.UserID, userAgent)
	if err != nil {
		utils.Fail(
			c,
			&utils.APIError{Code: http.StatusUnauthorized, Message: "Invalid or expired session"},
			err,
		)
		return
	}
	// 		- not expired
	if time.Now().After(session.ExpiresAt) {
		utils.Fail(
			c,
			&utils.APIError{Code: http.StatusUnauthorized, Message: "Invalid or expired session"},
			nil,
		)
		return
	}
	// 		- not revoked
	if session.Revoked {
		utils.Fail(
			c,
			&utils.APIError{Code: http.StatusUnauthorized, Message: "Invalid or expired session"},
			nil,
		)
		return
	}

	// validate refresh token
	if ok := utils.VerifyToken(session.RefreshToken, refreshToken, s.env.HashSecret); !ok {
		utils.Fail(
			c,
			&utils.APIError{Code: http.StatusUnauthorized, Message: "Invalid or expired session"},
			errors.New("couldn't verify the refresh token"),
		)
		return
	}

	// get user Role
	role, err := userRepo.GetRole(c, db, session.UserID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	userID := strconv.Itoa(int(session.UserID))

	// Rotate the refresh token
	newAccess, newRefresh, err := s.generateTokens(userID, string(role))
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	hashedRefresh, err := utils.HashToken(newRefresh, s.env.HashSecret)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	refreshExpTime := utils.GetExpTimeAfterDays(s.env.RefreshTokenExpInDays)

	session.RefreshToken = hashedRefresh
	session.ExpiresAt = refreshExpTime
	err = sessionRepo.Update(c, db, &session)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	// Return access and refresh tokens
	s.setRefreshCookie(c, newRefresh)

	utils.Success(c, refreshTokenRes{
		AccessToken: newAccess,
	})
}
