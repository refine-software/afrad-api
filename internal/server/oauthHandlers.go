package server

import (
	"errors"
	"slices"

	"github.com/gin-gonic/gin"
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

func (s *Server) googleCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		utils.Fail(c, utils.ErrUnauthorized, err)
		return
	}
	// fmt.Println("User ID:", user.UserID)
	// fmt.Println("User Avatar:", user.AvatarURL)
	// fmt.Println("User Email:", user.Email)
	// fmt.Println("User FirstName:", user.FirstName)
	// fmt.Println("User LastName:", user.LastName)
	// fmt.Println("User Name:", user.Name)
	// fmt.Println("Provider:", user.Provider)

	// Lookup user
	userRepo := s.db.User()
	db := s.db.Pool()

	u, err := userRepo.GetUserByEmail(c, db, user.Email)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	admins := []string{
		"ali93456@gmail.com",
		"bruhgg596@gmail.com",
	}
	userRole := models.RoleUser
	if slices.Contains(admins, user.Email) {
		userRole = models.RoleAdmin
	}

	// if user exists update it
	if u != nil {
		if user.FirstName == "" && user.Name != "" {
			u.FirstName = user.Name
		} else {
			u.FirstName = user.FirstName
		}
		u.LastName = user.LastName
		u.Image = user.AvatarURL
		u.Role = userRole

		err = userRepo.UpdateUser(c, db, u)
		if err != nil {
			apiErr := utils.MapDBErrorToAPIError(err, "user")
			utils.Fail(c, apiErr, err)
			return
		}
	} else {
		// if user doesn't exists create it
		u = &models.User{}
		if user.FirstName == "" && user.Name != "" {
			u.FirstName = user.Name
		} else {
			u.FirstName = user.FirstName
		}
		u.LastName = user.LastName
		u.Image = user.AvatarURL
		u.Role = userRole
		u.Email = user.Email

		err := userRepo.CreateUser(c, db, u)
		if err != nil {
			apiErr := utils.MapDBErrorToAPIError(err, "user")
			utils.Fail(c, apiErr, err)
			return
		}

		oauthRepo := s.db.Oauth()
		err = oauthRepo.Create(c, db, &models.OAuth{
			UserID:     u.ID,
			Provider:   user.Provider,
			ProviderID: user.UserID,
		})
		if err != nil {
			apiErr := utils.MapDBErrorToAPIError(err, "oauth")
			utils.Fail(c, apiErr, err)
			return
		}
	}

	// Generate access and refresh token
	// Create a session for the user
	// Use a transaction to do all of these db calls
}
