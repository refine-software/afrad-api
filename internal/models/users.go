package models

import (
	"time"
)

type Role string

type User struct {
	ID              int
	FirstName       string
	LastName        string
	Image           string
	PhoneNumber     string
	Email           string
	PasswordHash    string
	Role            Role
	IsPhoneVerified bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
