package models

import "time"

type Order struct {
	ID          int
	Town        string
	Street      string
	Address     string
	Name        string
	PhoneNumber string
	TotalPrice  int
	CitiesID    int
	UserID      int
	OrderStatus OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CancelledAt time.Time
}
