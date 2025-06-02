package models

import "time"

type Cart struct {
	ID         int
	TotalPrice int
	Quantity   int
	CreatedAt  time.Time
	UserID     int
}
