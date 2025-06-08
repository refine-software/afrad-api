package models

import "time"

type User struct {
	ID        int32     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Image     string    `json:"image"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type OAuth struct {
	UserID     int32
	Provider   string
	ProviderID string
}

type LocalAuth struct {
	UserID            int32
	PhoneNumber       string
	IsAccountVerified bool
	PasswordHash      string
}

type Session struct {
	ID           int32
	Revoked      bool
	UserAgent    string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserID       int32
}

type PasswordReset struct {
	ID        int32
	OtpCode   string
	IsUsed    bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UserID    int32
}

type PhoneVerification struct {
	ID        int32
	OtpCode   string
	IsUsed    bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UserID    int32
}
