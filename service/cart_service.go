package service

import (
	"context"

	"github.com/gunawan98/golang-restfull-api/model/web"
)

type CartService interface {
	AvailableCart(ctx context.Context, userId float64) []web.CartResponse
	CreateNewCart(ctx context.Context, request web.CartCreateRequest, userId float64) web.CartResponse
	AddProductToCart(ctx context.Context, userId float64, cartId int, request web.CartItemCreateRequest) web.CartItemResponse
	GetCartDetails(ctx context.Context, cartId int) (web.CartResponse, []web.CartItemResponse)
}
