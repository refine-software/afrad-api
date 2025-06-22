package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	ID              int32       `json:"id"`
	Name            string      `json:"name"`
	Details         pgtype.Text `json:"details"`
	Thumbnail       string      `json:"thumbnail"`
	CreatedAt       time.Time   `json:"-"`
	UpdatedAt       time.Time   `json:"-"`
	BrandID         int32       `json:"-"`
	ProductCategory int32       `json:"-"`
}

type ProductVariant struct {
	ID        int32     `json:"id"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	ProductID int32     `json:"-"`
	ColorID   int32     `json:"-"`
	SizeID    int32     `json:"-"`
}

type Category struct {
	ID       int32       `json:"id"`
	Name     string      `json:"name"`
	ParentID pgtype.Int4 `json:"parentId"`
}

type RatingReview struct {
	ID        int32       `json:"id"`
	Rating    int32       `json:"rating"`
	Review    pgtype.Text `json:"review"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	UserID    int32       `json:"-"`
	ProductID int32       `json:"-"`
}

type Image struct {
	ID          int32  `json:"id"`
	Image       string `json:"image"`
	LowResImage string `json:"lowResImage"`
	ProductID   int32  `json:"-"`
}

type Brand struct {
	ID    int32  `json:"id"`
	Brand string `json:"brand"`
}

type Size struct {
	ID    int32  `json:"id"`
	Size  string `json:"size"`
	Label string `json:"label"`
}

type Color struct {
	ID    int32  `json:"id"`
	Color string `json:"color"`
}
