package models

import (
	"time"
)

type Discount struct {
	ID            int
	DiscountType  string
	DiscountValue float64 // ??
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
