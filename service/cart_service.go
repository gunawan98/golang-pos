package service

import (
	"context"

	"github.com/gunawan98/golang-restfull-api/model/web"
)

type CartService interface {
	AvailableCart(ctx context.Context, userId float64) []web.CartResponse
	FinishedCart(ctx context.Context, userId float64, page int, perPage int) ([]web.CartResponse, int)
	CreateNewCart(ctx context.Context, userId float64) web.CartResponse
	AddProductToCart(ctx context.Context, userId float64, cartId int, request web.CartItemCreateRequest) web.CartItemResponse
	UpdateProductInCart(ctx context.Context, userId float64, cartId int, request web.CartItemUpdateRequest) web.CartItemResponse
	GetCartDetails(ctx context.Context, cartId int) (web.CartResponse, []web.CartItemWithProductResponse)
	DeleteCart(ctx context.Context, userId float64, cartId int)
	DeleteCartItem(ctx context.Context, userId float64, cartItemId int)
}
