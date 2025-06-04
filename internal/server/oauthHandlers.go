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
	db := s.db.Pool()

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
	sessExpTime := time.Now().Add((time.Hour * 24) * time.Duration(s.env.RefreshTokenExpInDays))
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

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int((time.Hour * 24 * time.Duration(s.env.RefreshTokenExpInDays)).Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, oauthRes{
		AccessToken: accessToken,
		User:        *u,
	})
}
