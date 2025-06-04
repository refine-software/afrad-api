package models

import "time"

type Order struct {
	ID          int32
	Town        string
	Street      string
	Address     string
	Name        string
	PhoneNumber string
	TotalPrice  int
	CitiesID    int32
	UserID      int32
	OrderStatus OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CancelledAt time.Time
}

type OrderDetails struct {
	ID         int32
	Quantity   int
	TotalPrice int
	ProductID  int32
	OrderID    int32
}

type City struct {
	ID   int32
	City string
}
