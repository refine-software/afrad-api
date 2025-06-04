package server

import (
	"slices"

	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/models"
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
