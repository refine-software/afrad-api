package models

import "time"

type Discount struct {
	ID            int32
	DiscountType  string
	DiscountValue float64 // ??
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type VariantDiscount struct {
	ID         int32
	DiscountID int32
	VariantID  int32
}
