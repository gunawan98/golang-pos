package domain

import "time"

type Purchase struct {
	Id, CartID, CashierID, TotalAmount int
	PaymentMethod                      string
	CreatedAt                          time.Time
}
