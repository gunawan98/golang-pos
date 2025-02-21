package service

import (
	"context"

	"github.com/gunawan98/golang-restfull-api/model/web"
)

type ProductImageService interface {
	AddImage(ctx context.Context, request web.ProductImageCreateRequest) web.ProductImageResponse
	DeleteImage(ctx context.Context, imageId int)
	FindByProductId(ctx context.Context, productId int) []web.ProductImageResponse
	FindById(ctx context.Context, imageId int) web.ProductImageResponse
}
