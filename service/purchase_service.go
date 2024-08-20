package service

import (
	"context"

	"github.com/gunawan98/golang-restfull-api/model/web"
)

type PurchaseService interface {
	ConfirmPayment(ctx context.Context, request web.PurchaseCreateRequest) web.PurchaseResponse
}
