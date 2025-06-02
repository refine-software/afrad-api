package models

import "time"

type ProductVariant struct {
	ID        int
	Quantity  int
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
	ProductID int
	ColorID   int
	SizeID    int
}
