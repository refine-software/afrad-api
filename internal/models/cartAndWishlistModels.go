package models

import "time"

type Cart struct {
	ID         int32     `json:"id"`
	TotalPrice int       `json:"totalPrice,omitempty"`
	Quantity   int       `json:"quantity,omitempty"`
	CreatedAt  time.Time `json:"-"`
	UserID     int32     `json:"userId,omitempty"`
}

type CartItem struct {
	ID         int32
	Quantity   int
	TotalPrice int
	CreatedAt  time.Time
	CartID     int32
	ProductID  int32
}

type Wishlist struct {
	ID        int32
	UserID    int32
	ProductID int32
	CreatedAt time.Time
}
