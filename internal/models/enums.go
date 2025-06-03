package models

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

type OrderStatus string

const (
	OrderPlaced OrderStatus = "order_placed"
	InProgress  OrderStatus = "in_progress"
	Shipped     OrderStatus = "shipped"
	Delivered   OrderStatus = "delivered"
	Cancelled   OrderStatus = "cancelled"
)

func (o OrderStatus) IsValid() bool {
	switch o {
	case OrderPlaced, InProgress, Shipped, Delivered, Cancelled:
		return true
	}
	return false
}
