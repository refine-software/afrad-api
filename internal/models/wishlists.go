package models

import "time"

type Wishlist struct {
	ID        int
	UserID    int
	ProductID int
	CreatedAt time.Time
}
