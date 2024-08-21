package web

import "time"

type PurchaseCreateRequest struct {
	CartID        int       `validate:"required" json:"cart_id"`
	CashierID     int       `validate:"required" json:"cashier_id"`
	TotalAmount   int       `validate:"required" json:"total_amount"`
	PaymentMethod string    `validate:"oneof=cash credit-card ewallet other" json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}

type PurchaseResponse struct {
	CartID        int       `json:"cart_id"`
	CashierID     int       `json:"cashier_id"`
	TotalAmount   int       `json:"total_amount"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}
