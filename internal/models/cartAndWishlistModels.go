package models

import "time"

type Cart struct {
	ID         int32
	TotalPrice int
	Quantity   int
	CreatedAt  time.Time
	UserID     int32
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
