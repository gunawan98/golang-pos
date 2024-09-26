package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type PurchaseRepositoryImpl struct{}

func NewPurchaseRepository() PurchaseRepository {
	return &PurchaseRepositoryImpl{}
}

func (repository *PurchaseRepositoryImpl) AddPurchase(ctx context.Context, tx *sql.Tx, purchase domain.Purchase) domain.Purchase {
	SQL := "INSERT INTO purchase(cart_id, cashier_id, total_amount, paid, cash_back, payment_method, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, purchase.CartID, purchase.CashierID, purchase.TotalAmount, purchase.Paid, purchase.CashBack, purchase.PaymentMethod, purchase.CreatedAt)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	purchase.Id = int(id)
	return purchase
}

func (repository *PurchaseRepositoryImpl) GetPurchaseByCartId(ctx context.Context, tx *sql.Tx, userId int, cartId int) (domain.Purchase, error) {
	SQL := "SELECT id, cart_id, cashier_id, total_amount, paid, cash_back, payment_method, created_at FROM purchase WHERE cart_id=? AND cashier_id=?"
	rows, err := tx.QueryContext(ctx, SQL, cartId, userId)
	helper.PanicIfError(err)
	defer rows.Close()

	purchase := domain.Purchase{}
	if rows.Next() {
		err := rows.Scan(&purchase.Id, &purchase.CartID, &purchase.CashierID, &purchase.TotalAmount, &purchase.Paid, &purchase.CashBack, &purchase.PaymentMethod, &purchase.CreatedAt)
		helper.PanicIfError(err)
		return purchase, nil
	} else {
		return purchase, errors.New("purchase is not found")
	}
}
