package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	ID              int32
	Name            string
	Details         pgtype.Text
	Thumbnail       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BrandID         int32
	ProductCategory int32
}

type ProductVariant struct {
	ID        int32
	Quantity  int
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
	ProductID int32
	ColorID   int32
	SizeID    int32
}

type Category struct {
	ID       int32
	Name     string
	ParentID pgtype.Int4
}

type RatingReview struct {
	ID        int32
	Rating    int32
	Review    pgtype.Text
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int32
	ProductID int32
}

type Image struct {
	ID          int32
	Image       string
	LowResImage string
	ProductID   int32
}

type Brand struct {
	ID    int32
	Brand string
}

type Size struct {
	ID    int32
	Size  string
	Label string
}

type Color struct {
	ID    int32
	Color string
}
