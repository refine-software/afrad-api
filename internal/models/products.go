package models

import "time"

type Product struct {
	ID              int
	Name            string
	Details         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BrandID         int
	ProductCategory int
}
