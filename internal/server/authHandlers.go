package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
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
	fmt.Println("User ID:", user.UserID)
	fmt.Println("User Avatar:", user.AvatarURL)
	fmt.Println("User Email:", user.Email)
	fmt.Println("User FirstName:", user.FirstName)
	fmt.Println("User LastName:", user.LastName)
	fmt.Println("User Name:", user.Name)
	fmt.Println("Provider:", user.Provider)

	// Lookup user

	// if user exists update it

	// if user doesn't exists create it
}
