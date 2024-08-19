package repository

import (
	"context"
	"database/sql"

	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type ProductRepository interface {
	Save(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product
	Update(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product
	Delete(ctx context.Context, tx *sql.Tx, productId int)
	FindById(ctx context.Context, tx *sql.Tx, productId int) (domain.Product, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Product
	FindByBarcode(ctx context.Context, tx *sql.Tx, id int, barcode string) (domain.Product, error)
}
