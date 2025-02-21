package repository

import (
	"context"
	"database/sql"

	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type ProductImageRepository interface {
	Save(ctx context.Context, tx *sql.Tx, image domain.ProductImage) domain.ProductImage
	Delete(ctx context.Context, tx *sql.Tx, imageId int)
	FindByProductId(ctx context.Context, tx *sql.Tx, productId int) []domain.ProductImage
	FindById(ctx context.Context, tx *sql.Tx, imageId int) (domain.ProductImage, error)
}
