package models

import "time"

type PasswordReset struct {
	ID        int
	OtpCode   string
	IsUsed    bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UserID    int
}
