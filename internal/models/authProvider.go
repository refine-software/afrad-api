package models

import "github.com/google/uuid"

type AuthProvider struct {
	ID         int
	Provider   string
	ProviderID uuid.UUID
	UserID     int
}
