package models

import (
	"time"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Image     string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
