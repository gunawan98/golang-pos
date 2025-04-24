package web

import "time"

type CartResponse struct {
	Id        int       `json:"id"`
	CashierID int       `json:"cashier_id"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type CartDetailResponse struct {
	Id            int       `json:"id"`
	CashierID     int       `json:"cashier_id"`
	Completed     bool      `json:"completed"`
	CreatedAt     time.Time `json:"created_at"`
	TotalPurchase int       `json:"total_purchase"`
}

type CartItemCreateRequest struct {
	Barcode  string `validate:"required" json:"barcode"`
	Quantity int    `validate:"required,numeric,min=1" json:"quantity"`
}

type CartItemUpdateRequest struct {
	ProductID int `validate:"required,numeric,min=1" json:"product_id"`
	Quantity  int `validate:"required,numeric,min=1" json:"quantity"`
}

type CartItemResponse struct {
	Id         int `json:"id"`
	CartID     int `json:"cart_id"`
	ProductID  int `json:"product_id"`
	Quantity   int `json:"quantity"`
	UnitPrice  int `json:"unit_price"`
	TotalPrice int `json:"total_price"`
}

type CartItemWithProductResponse struct {
	Id           int    `json:"id"`
	CartID       int    `json:"cart_id"`
	ProductID    int    `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`
	Quantity     int    `json:"quantity"`
	UnitPrice    int    `json:"unit_price"`
	TotalPrice   int    `json:"total_price"`
}
