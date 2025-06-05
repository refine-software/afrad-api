package server

import (
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

func getUserRole(email string) models.Role {
	admins := []string{
		"ali93456@gmail.com",
		"bruhgg596@gmail.com",
	}
	if slices.Contains(admins, email) {
		return models.RoleAdmin
	}
	return models.RoleUser
}

func getNameFallback(firstName, name string) string {
	if firstName == "" && name != "" {
		return name
	}
	return firstName
}

func (s *Server) generateTokens(userID, role string) (access, refresh string, err error) {
	access, err = auth.GenerateAccessToken(
		userID,
		role,
		s.env.AccessTokenSecret,
		s.env.AccessTokenExpInMin,
	)
	if err != nil {
		return
	}
	refresh, err = auth.GenerateRefreshToken(
		userID,
		s.env.RefreshTokenSecret,
		s.env.RefreshTokenExpInDays,
	)
	return
}

func getHeader(c *gin.Context, key string) string {
	header := strings.TrimSpace(c.GetHeader(key))
	if header == "" {
		utils.Fail(
			c,
			utils.ErrHeaderMissing(key),
			nil,
		)
		return ""
	}
	return header
}

func (s *Server) setCookie(c *gin.Context, cookieName, cookieVal string) {
	var secure bool
	if s.env.Environment == "prod" {
		secure = true
	}

	expTimeInSec := int((time.Hour * 24 * time.Duration(s.env.RefreshTokenExpInDays)).Seconds())

	c.SetCookie(
		cookieName,
		cookieVal,
		expTimeInSec,
		"/",
		"",
		secure,
		true,
	)
}

func (s *Server) setRefreshCookie(c *gin.Context, refreshToken string) {
	s.setCookie(c, "refreshToken", refreshToken)
}

func getExpTimeAfterDays(numOfDays int) time.Time {
	return time.Now().Add((time.Hour * 24) * time.Duration(numOfDays))
}

func getExpTimeAfterMins(numOfMins int) time.Time {
	return time.Now().Add((time.Minute * time.Duration(numOfMins)))
}
