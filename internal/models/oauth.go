package models

import "github.com/google/uuid"

type OAuth struct {
	UserID     int
	Provider   string
	ProviderID uuid.UUID
}
