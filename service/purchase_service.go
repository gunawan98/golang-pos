package service

import (
	"context"

	"github.com/gunawan98/golang-restfull-api/model/web"
)

type PurchaseService interface {
	ConfirmPayment(ctx context.Context, request web.PurchaseCreateRequest, userId float64) (web.PurchaseResponse, []web.CartItemWithProductResponse)
	GetFinishedPayment(ctx context.Context, userId float64, cartId int) (web.PurchaseResponse, []web.CartItemWithProductResponse)
}
