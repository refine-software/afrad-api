package models

type OrderDetails struct {
	ID         int
	Quantity   int
	TotalPrice int
	ProductID  int
	OrderID    int
}
