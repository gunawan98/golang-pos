package service

import (
	"context"

	"github.com/gunawan98/golang-restfull-api/model/web"
)

type CartService interface {
	CreateNewCart(ctx context.Context, request web.CartCreateRequest, userId float64) web.CartResponse
	AddProductToCart(ctx context.Context, cartId int, request web.CartItemCreateRequest) web.CartItemResponse
	GetCartDetails(ctx context.Context, cartId int) (web.CartResponse, []web.CartItemResponse)
}
