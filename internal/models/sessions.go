package models

import "time"

type Session struct {
	ID           int
	Revoked      bool
	UserAgent    string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserID       int
}
