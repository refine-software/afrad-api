package server

import (
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

	userRepo := s.db.User()
	localAuthRepo := s.db.LocalAuth()
	sessionRepo := s.db.Session()
	db := s.db.Pool()

	dbErr := userRepo.CheckEmailExistence(ctx, db, req.Email)
	if dbErr != nil {
		utils.Fail(ctx, utils.ErrInvalidCredentials, dbErr)
		return
	}

	user, dbErr := userRepo.GetByEmail(ctx, db, req.Email)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "user")
		utils.Fail(ctx, apiErr, dbErr)
		return
	}

	localAuth, dbErr := localAuthRepo.Get(ctx, db, user.ID)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "user")
		utils.Fail(ctx, apiErr, dbErr)
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

	hashedNewRefreshToken, err := utils.HashToken(newRefreshToken, s.env.HashSecret)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	// create or update session
	var session models.Session
	session, dbErr = sessionRepo.GetByUserIDAndUserAgent(
		ctx,
		db,
		user.ID,
		ctx.Request.UserAgent(),
	)
	if dbErr != nil && dbErr.Message != database.ErrNotFound {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "session")
		utils.Fail(ctx, apiErr, dbErr)
		return
	}
	sessExpTime := utils.GetExpTimeAfterDays(s.env.RefreshTokenExpInDays)
	if dbErr != nil && dbErr.Message == database.ErrNotFound {
		session = models.Session{
			UserID:       user.ID,
			RefreshToken: hashedNewRefreshToken,
			ExpiresAt:    sessExpTime,
			UserAgent:    ctx.Request.UserAgent(),
		}

		dbErr = sessionRepo.Create(ctx, db, &session)
		if dbErr != nil {
			apiErr := utils.MapDBErrorToAPIError(dbErr, "session")
			utils.Fail(ctx, apiErr, dbErr)
			return
		}
	} else {
		session.Revoked = false
		session.RefreshToken = hashedNewRefreshToken
		session.ExpiresAt = sessExpTime
		dbErr = sessionRepo.Update(ctx, db, &session)
		if dbErr != nil {
			apiErr := utils.MapDBErrorToAPIError(dbErr, "session")
			utils.Fail(ctx, apiErr, dbErr)
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

	db := s.db.Pool()
	sessionRepo := s.db.Session()

	userAgent := getHeader(c, "User-Agent")
	if userAgent == "" {
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	session, dbErr := sessionRepo.GetByUserIDAndUserAgent(c, db, int32(userID), userAgent)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "session")
		utils.Fail(c, apiErr, dbErr)
		return
	}

	// Check refresh token validity
	ok := utils.VerifyToken(session.RefreshToken, refreshToken, s.env.RefreshTokenSecret)
	if !ok {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	// revoke the session
	session.Revoked = true
	dbErr = sessionRepo.Update(c, db, &session)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "session")
		utils.Fail(c, apiErr, dbErr)
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

	db := s.db.Pool()
	sessionRepo := s.db.Session()

	dbErr := sessionRepo.RevokeAllOfUser(c, db, int32(userID))
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "session")
		utils.Fail(c, apiErr, dbErr)
		return
	}

	utils.NoContent(c)
}
