package web

import "time"

type CartCreateRequest struct {
	CashierID int  `validate:"required" json:"cashier_id"`
	Completed bool `json:"completed"`
}

type CartResponse struct {
	Id        int       `json:"id"`
	CashierID int       `json:"cashier_id"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type CartItemCreateRequest struct {
	ProductID  int `validate:"required" json:"product_id"`
	Quantity   int `validate:"required" json:"quantity"`
	UnitPrice  int `validate:"required" json:"unit_price"`
	TotalPrice int `validate:"required" json:"total_price"`
}

type CartItemResponse struct {
	Id         int `json:"id"`
	CartID     int `json:"cart_id"`
	ProductID  int `json:"product_id"`
	Quantity   int `json:"quantity"`
	UnitPrice  int `json:"unit_price"`
	TotalPrice int `json:"total_price"`
}
