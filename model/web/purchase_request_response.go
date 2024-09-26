package web

import "time"

type PurchaseCreateRequest struct {
	CartID        int       `validate:"required" json:"cart_id"`
	Paid          int       `validate:"required" json:"paid"`
	PaymentMethod string    `validate:"oneof=cash credit-card ewallet other" json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}

type PurchaseResponse struct {
	CartID        int       `json:"cart_id"`
	CashierID     int       `json:"cashier_id"`
	TotalAmount   int       `json:"total_amount"`
	Paid          int       `json:"paid"`
	CashBack      int       `json:"cash_back"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}
