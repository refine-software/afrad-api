package models

type LocalAuth struct {
	UserID          int
	PhoneNumber     string
	IsPhoneVerified bool
	PasswordHash    string
}
