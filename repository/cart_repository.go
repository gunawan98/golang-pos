package repository

import (
	"context"
	"database/sql"

	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type CartRepository interface {
	CreateCart(ctx context.Context, tx *sql.Tx, cart domain.Cart) domain.Cart
	GetCartById(ctx context.Context, tx *sql.Tx, cartId int) (domain.Cart, error)
	AddItemToCart(ctx context.Context, tx *sql.Tx, cartItem domain.CartItem) domain.CartItem
	FindItemByCartAndProduct(ctx context.Context, tx *sql.Tx, cartId int, productId int) (domain.CartItem, error)
	UpdateCartItem(ctx context.Context, tx *sql.Tx, cartItem domain.CartItem) error
	GetItemsByCartId(ctx context.Context, tx *sql.Tx, cartId int) []domain.CartItem
	UpdateCartStatus(ctx context.Context, tx *sql.Tx, cartId int) error
	FindAvailableCart(ctx context.Context, tx *sql.Tx, userId int) []domain.Cart
}
