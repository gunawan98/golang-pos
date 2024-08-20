package repository

import (
	"context"
	"database/sql"

	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type PurchaseRepositoryImpl struct{}

func NewPurchaseRepository() PurchaseRepository {
	return &PurchaseRepositoryImpl{}
}

func (repository *PurchaseRepositoryImpl) AddPurchase(ctx context.Context, tx *sql.Tx, purchase domain.Purchase) domain.Purchase {
	SQL := "INSERT INTO purchase(cart_id, cashier_id, total_amount, payment_method, created_at) VALUES (?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, purchase.CartID, purchase.CashierID, purchase.TotalAmount, purchase.PaymentMethod, purchase.CreatedAt)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	purchase.Id = int(id)
	return purchase
}
