package domain

import "time"

type Cart struct {
	Id, CashierID int
	Completed     bool
	CreatedAt     time.Time
}
