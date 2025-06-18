package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID          int32       `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    pgtype.Text `json:"lastName"`
	Image       pgtype.Text `json:"image"`
	Email       string      `json:"email"`
	PhoneNumber pgtype.Text `json:"phoneNumber"`
	Role        Role        `json:"role"`
	CreatedAt   time.Time   `json:"-"`
	UpdatedAt   time.Time   `json:"-"`
}

type OAuth struct {
	UserID     int32
	Provider   string
	ProviderID string
}

type LocalAuth struct {
	UserID            int32
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

type AccountVerificationCode struct {
	ID        int32
	OtpCode   string
	IsUsed    bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UserID    int32
}
