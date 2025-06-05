package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

func (s *Server) loginWithGoogle(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

type upsertResult struct {
	User     *models.User
	IsNew    bool
	APIError *utils.APIError
	Err      error
}

func (s *Server) upsertUser(
	c *gin.Context,
	db database.Querier,
	user goth.User,
	role models.Role,
) (u *models.User, resultErr upsertResult) {
	userRepo := s.db.User()
	oauthRepo := s.db.Oauth()

	u, err := userRepo.Get(c, db, user.Email)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, upsertResult{
			APIError: utils.MapDBErrorToAPIError(err, "user"),
			Err:      err,
		}
	}

	if u != nil {
		// update user
		u.FirstName = getNameFallback(user.FirstName, user.Name)
		u.LastName = user.LastName
		u.Image = user.AvatarURL
		u.Role = role

		err = userRepo.Update(c, db, u)
		if err != nil {
			return nil, upsertResult{
				APIError: utils.MapDBErrorToAPIError(err, "user"),
				Err:      err,
			}
		}
		return u, upsertResult{}
	}

	// create user
	u = &models.User{
		FirstName: getNameFallback(user.FirstName, user.Name),
		LastName:  user.LastName,
		Image:     user.AvatarURL,
		Email:     user.Email,
		Role:      role,
	}
	userID, err := userRepo.Create(c, db, u)
	if err != nil {
		return nil, upsertResult{APIError: utils.MapDBErrorToAPIError(err, "user"), Err: err}
	}
	u.ID = int32(userID)

	err = oauthRepo.Create(c, db, &models.OAuth{
		UserID:     u.ID,
		Provider:   user.Provider,
		ProviderID: user.UserID,
	})
	if err != nil {
		return nil, upsertResult{
			APIError: utils.MapDBErrorToAPIError(err, "oauth"),
			Err:      err,
		}
	}

	return u, upsertResult{}
}

type oauthRes struct {
	AccessToken string      `json:"accessToken"`
	User        models.User `json:"user"`
}

func (s *Server) googleCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		utils.Fail(c, utils.ErrUnauthorized, err)
		return
	}

	sessionRepo := s.db.Session()
	db, err := s.db.BeginTx(c)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = db.Rollback(c)
			panic(p)
		}
	}()

	u, upsertErr := s.upsertUser(c, db, user, getUserRole(user.Email))
	if upsertErr.Err != nil || u == nil {
		utils.Fail(c, upsertErr.APIError, upsertErr.Err)
		return
	}

	userIDStr := strconv.Itoa(int(u.ID))
	accessToken, refreshToken, err := s.generateTokens(userIDStr, string(u.Role))
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	// hash refresh token
	hashedRefresh, err := utils.HashToken(refreshToken, s.env.HashSecret)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	var session models.Session
	session, err = sessionRepo.GetByUserIDAndUserAgent(c, db, u.ID, c.Request.UserAgent())
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		apiErr := utils.MapDBErrorToAPIError(err, "session")
		utils.Fail(c, apiErr, err)
		return
	}
	sessExpTime := getExpTimeAfterDays(s.env.RefreshTokenExpInDays)
	if errors.Is(err, database.ErrNotFound) {
		session = models.Session{
			UserID:       u.ID,
			RefreshToken: hashedRefresh,
			ExpiresAt:    sessExpTime,
			UserAgent:    c.Request.UserAgent(),
		}

		err = sessionRepo.Create(c, db, &session)
		if err != nil {
			apiErr := utils.MapDBErrorToAPIError(err, "session")
			utils.Fail(c, apiErr, err)
			return
		}
	} else {
		session.Revoked = false
		session.RefreshToken = hashedRefresh
		session.ExpiresAt = sessExpTime
		err = sessionRepo.Update(c, db, &session)
		if err != nil {
			apiErr := utils.MapDBErrorToAPIError(err, "session")
			utils.Fail(c, apiErr, err)
			return
		}
	}

	err = db.Commit(c)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	s.setRefreshCookie(c, refreshToken)

	utils.Success(c, "You have loged in successfully", oauthRes{
		AccessToken: accessToken,
		User:        *u,
	})
}

type refreshTokenReq struct {
	UserID int32 `json:"userId"`
}

type refreshTokenRes struct {
	AccessToken string `json:"accessToken"`
}

func (s *Server) refreshTokenOauth(c *gin.Context) {
	var req refreshTokenReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(c, utils.ErrBadRequest, nil)
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

	userAgent := c.Request.UserAgent()
	if userAgent == "" {
		utils.Fail(
			c,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "User-Agent header is required",
			},
			err,
		)
		return
	}

	sessionRepo := s.db.Session()
	userRepo := s.db.User()
	db, err := s.db.BeginTx(c)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = db.Rollback(c)
			panic(p)
		}
	}()

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
	if ok := utils.VerifyToken(session.RefreshToken, refreshToken, s.env.RefreshTokenSecret); !ok {
		utils.Fail(
			c,
			&utils.APIError{Code: http.StatusUnauthorized, Message: "Invalid or expired session"},
			nil,
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
	access, refresh, err := s.generateTokens(userID, string(role))
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	hashedRefresh, err := utils.HashToken(refresh, s.env.RefreshTokenSecret)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	refreshExpTime := getExpTimeAfterDays(s.env.RefreshTokenExpInDays)

	session.RefreshToken = hashedRefresh
	session.ExpiresAt = refreshExpTime
	err = sessionRepo.Update(c, db, &session)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	err = db.Commit(c)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	// Return access and refresh tokens
	s.setRefreshCookie(c, refreshToken)

	utils.Success(c, "your tokens have been refreshed", refreshTokenRes{
		AccessToken: access,
	})
}
