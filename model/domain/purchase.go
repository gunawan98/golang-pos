package domain

import "time"

type Purchase struct {
	Id, CartID, CashierID, TotalAmount, Paid, CashBack int
	PaymentMethod                                      string
	CreatedAt                                          time.Time
}
