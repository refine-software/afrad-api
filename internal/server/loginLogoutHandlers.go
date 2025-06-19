package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

type loginReq struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Email/Password Login
// @Description  Logs in a user using email and password. Returns an access token and user data.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginReq  body      loginReq       true  "Login request"
// @Success      200       {object}  loginResDocs       "Successful login with access token and user info"
// @Failure      400       {object}  utils.APIError "Invalid request body"
// @Failure      401       {object}  utils.APIError "Invalid credentials or unverified account"
// @Failure      500       {object}  utils.APIError "Internal server error"
// @Router       /auth/login [post]
func (s *Server) login(ctx *gin.Context) {
	var req loginReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	userRepo := s.DB.User()
	localAuthRepo := s.DB.LocalAuth()
	sessionRepo := s.DB.Session()
	db := s.DB.Pool()

	fmt.Println("CheckEmailExistence")
	err = userRepo.CheckEmailExistence(ctx, db, req.Email)
	if err != nil {
		utils.Fail(ctx, utils.ErrInvalidCredentials, err)
		return
	}

	fmt.Println("GetByEmail")
	user, err := userRepo.GetByEmail(ctx, db, req.Email)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(ctx, apiErr, err)
		return
	}

	fmt.Println("Get")
	localAuth, err := localAuthRepo.Get(ctx, db, user.ID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(ctx, apiErr, err)
		return
	}

	// check password
	if err = utils.VerifyPassword(localAuth.PasswordHash, req.Password); err != nil {
		utils.Fail(ctx, utils.ErrInvalidCredentials, err)
		return
	}

	// check if user is verified
	if !localAuth.IsAccountVerified {
		utils.Fail(ctx, &utils.APIError{
			Code:    http.StatusUnauthorized,
			Message: "your account isn't verified yet",
		}, nil)
		return
	}
	userIDStr := strconv.Itoa(int(user.ID))
	newAccessToken, newRefreshToken, err := s.generateTokens(userIDStr, string(user.Role))
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	hashedNewRefreshToken, err := utils.HashToken(newRefreshToken, s.Env.HashSecret)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	fmt.Println("GetByUserIDAndUserAgent")
	// create or update session
	var session models.Session
	session, err = sessionRepo.GetByUserIDAndUserAgent(
		ctx,
		db,
		user.ID,
		ctx.Request.UserAgent(),
	)
	if err != nil && database.IsDBNotFoundErr(err) {
		apiErr := utils.MapDBErrorToAPIError(err, "session")
		utils.Fail(ctx, apiErr, err)
		return
	}
	sessExpTime := utils.GetExpTimeAfterDays(s.Env.RefreshTokenExpInDays)
	if err != nil && database.IsDBNotFoundErr(err) {
		session = models.Session{
			UserID:       user.ID,
			RefreshToken: hashedNewRefreshToken,
			ExpiresAt:    sessExpTime,
			UserAgent:    ctx.Request.UserAgent(),
		}

		fmt.Println("Create")
		err = sessionRepo.Create(ctx, db, &session)
		if err != nil {
			apiErr := utils.MapDBErrorToAPIError(err, "session")
			utils.Fail(ctx, apiErr, err)
			return
		}
	} else {
		session.Revoked = false
		session.RefreshToken = hashedNewRefreshToken
		session.ExpiresAt = sessExpTime
		fmt.Println("Update")
		err = sessionRepo.Update(ctx, db, &session)
		if err != nil {
			apiErr := utils.MapDBErrorToAPIError(err, "session")
			utils.Fail(ctx, apiErr, err)
			return
		}
	}

	s.setRefreshCookie(ctx, newRefreshToken)

	utils.Success(ctx, loginRes{
		AccessToken: newAccessToken,
		User:        *user,
	})
}

// @Summary      Logout
// @Description  Logs out the currently authenticated user by revoking the session and clearing the refresh token cookie.
// @Tags         User
// @Security     BearerAuth
// @Produce      json
// @Success      204  "Successfully logged out"
// @Failure      400  {object}  utils.APIError  "Missing refresh token or invalid request"
// @Failure      401  {object}  utils.APIError  "Unauthorized or invalid session"
// @Failure      500  {object}  utils.APIError  "Internal server error"
// @Router       /user/logout [post]
func (s *Server) logout(c *gin.Context) {
	claims := auth.GetAccessClaims(c)
	if claims == nil {
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.Fail(
			c,
			&utils.APIError{Code: http.StatusBadRequest, Message: "Missing refresh token cookie"},
			err,
		)
		return
	}

	db := s.DB.Pool()
	sessionRepo := s.DB.Session()

	userAgent := getHeader(c, "User-Agent")
	if userAgent == "" {
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	session, err := sessionRepo.GetByUserIDAndUserAgent(c, db, int32(userID), userAgent)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "session")
		utils.Fail(c, apiErr, err)
		return
	}

	// Check refresh token validity
	ok := utils.VerifyToken(session.RefreshToken, refreshToken, s.Env.RefreshTokenSecret)
	if !ok {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	// revoke the session
	session.Revoked = true
	err = sessionRepo.Update(c, db, &session)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "session")
		utils.Fail(c, apiErr, err)
		return
	}

	utils.NoContent(c)
}

// @Summary      Logout from All Sessions
// @Description  Revokes all active sessions for the authenticated user across all devices.
// @Tags         User
// @Security     BearerAuth
// @Produce      json
// @Success      204  "Successfully logged out from all sessions"
// @Failure      401  {object}  utils.APIError  "Unauthorized or invalid token"
// @Failure      500  {object}  utils.APIError  "Internal server error"
// @Router       /user/logout/all [post]
func (s *Server) logoutFromAllSessions(c *gin.Context) {
	claims := auth.GetAccessClaims(c)
	if claims == nil {
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	db := s.DB.Pool()
	sessionRepo := s.DB.Session()

	err = sessionRepo.RevokeAllOfUser(c, db, int32(userID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "session")
		utils.Fail(c, apiErr, err)
		return
	}

	utils.NoContent(c)
}
