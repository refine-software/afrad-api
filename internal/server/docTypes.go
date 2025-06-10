package server

import "github.com/refine-software/afrad-api/internal/models"

type userDocs struct {
	ID          int32       `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	Image       string      `json:"image"`
	Email       string      `json:"email"`
	PhoneNumber string      `json:"PhoneNumber"`
	Role        models.Role `json:"role"`
}

type loginResDocs struct {
	AccessToken string   `json:"accessToken"`
	User        userDocs `json:"user"`
}
