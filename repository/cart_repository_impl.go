package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type CartRepositoryImpl struct {
}

func NewCartRepository() CartRepository {
	return &CartRepositoryImpl{}
}

func (r *CartRepositoryImpl) CreateCart(ctx context.Context, tx *sql.Tx, cart domain.Cart) domain.Cart {
	SQL := "INSERT INTO cart (cashier_id, completed) VALUES (?, ?)"
	result, err := tx.ExecContext(ctx, SQL, cart.CashierID, cart.Completed)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	cart.Id = int(id)
	return cart
}

func (r *CartRepositoryImpl) GetCartById(ctx context.Context, tx *sql.Tx, cartId int) (domain.Cart, error) {
	SQL := "SELECT id, cashier_id, completed, created_at FROM cart WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, cartId)
	helper.PanicIfError(err)
	defer rows.Close()

	cart := domain.Cart{}
	if rows.Next() {
		err := rows.Scan(&cart.Id, &cart.CashierID, &cart.Completed, &cart.CreatedAt)
		helper.PanicIfError(err)
		return cart, nil
	} else {
		return cart, errors.New("cart is not found")
	}
}

func (r *CartRepositoryImpl) AddItemToCart(ctx context.Context, tx *sql.Tx, cartItem domain.CartItem) domain.CartItem {
	SQL := "INSERT INTO cart_item (cart_id, product_id, quantity, unit_price, total_price) VALUES (?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, cartItem.CartID, cartItem.ProductID, cartItem.Quantity, cartItem.UnitPrice, cartItem.TotalPrice)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	cartItem.Id = int(id)
	return cartItem
}

func (r *CartRepositoryImpl) GetItemsByCartId(ctx context.Context, tx *sql.Tx, cartId int) []domain.CartItem {
	SQL := "SELECT id, cart_id, product_id, quantity, unit_price, total_price FROM cart_item WHERE cart_id = ?"
	rows, err := tx.QueryContext(ctx, SQL, cartId)
	helper.PanicIfError(err)
	defer rows.Close()

	var items []domain.CartItem
	for rows.Next() {
		item := domain.CartItem{}
		err := rows.Scan(&item.Id, &item.CartID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.TotalPrice)
		helper.PanicIfError(err)
		items = append(items, item)
	}

	return items
}
