package models

import "time"

type CartItem struct {
	ID         int
	Quantity   int
	TotalPrice int
	CreatedAt  time.Time
	CartID     int
	ProductID  int
}
