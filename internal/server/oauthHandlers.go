package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

// @Summary      Start Google OAuth Login
// @Description  Redirects the user to Google's OAuth 2.0 login screen.
// @Tags         OAuth
// @Accept       json
// @Produce      json
// @Success      302  {string}  string  "Redirect to Google"
// @Failure      500  {object}  utils.APIError  "Internal Server Error"
// @Router       /oauth/google/login [get]
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

	u, err := userRepo.GetByEmail(c, db, user.Email)
	if err != nil && database.IsDBNotFoundErr(err) {
		return nil, upsertResult{
			APIError: utils.MapDBErrorToAPIError(err, "user"),
			Err:      err,
		}
	}

	if u != nil {
		// update user
		u.FirstName = getNameFallback(user.FirstName, user.Name)
		u.LastName = pgtype.Text{String: user.LastName, Valid: true}
		u.Image = pgtype.Text{String: user.AvatarURL, Valid: user.AvatarURL != ""}
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
		FirstName:   getNameFallback(user.FirstName, user.Name),
		LastName:    pgtype.Text{String: user.LastName, Valid: true},
		Image:       pgtype.Text{String: user.AvatarURL, Valid: user.AvatarURL != ""},
		Email:       user.Email,
		PhoneNumber: pgtype.Text{},
		Role:        role,
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

type loginRes struct {
	AccessToken string      `json:"accessToken"`
	User        models.User `json:"user"`
}

// @Summary      Google OAuth Callback
// @Description  Handles the Google OAuth callback, authenticates the user, and returns a JWT access token.
// @Tags         OAuth
// @Accept       json
// @Produce      json
// @Param        code   query     string  true  "OAuth authorization code"
// @Param        state  query     string  false "OAuth state (if used)"
// @Success      200    {object}  loginResDocs   "Successful login with JWT token and user data"
// @Failure      400    {object}  utils.APIError  "Bad request or invalid input"
// @Failure      401    {object}  utils.APIError  "Unauthorized - Invalid OAuth token"
// @Failure      500    {object}  utils.APIError  "Internal Server Error"
// @Router       /oauth/google/callback [get]
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

	committed := false
	defer func() {
		if p := recover(); p != nil {
			_ = db.Rollback(c)
			panic(p)
		} else if !committed {
			_ = db.Rollback(c)
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
	if err != nil && database.IsDBNotFoundErr(err) {
		apiErr := utils.MapDBErrorToAPIError(err, "session")
		utils.Fail(c, apiErr, err)
		return
	}
	sessExpTime := utils.GetExpTimeAfterDays(s.env.RefreshTokenExpInDays)
	if err != nil && database.IsDBNotFoundErr(err) {
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
	committed = true

	s.setRefreshCookie(c, refreshToken)

	utils.Success(c, loginRes{
		AccessToken: accessToken,
		User:        *u,
	})
}
