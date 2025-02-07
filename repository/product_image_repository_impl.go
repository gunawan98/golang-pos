package repository

import (
	"context"
	"database/sql"

	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type ProductImageRepositoryImpl struct {
}

func NewProductImageRepository() ProductImageRepository {
	return &ProductImageRepositoryImpl{}
}

func (repository *ProductImageRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, image domain.ProductImage) domain.ProductImage {
	SQL := "INSERT INTO product_image(product_id, url) VALUES (?, ?)"
	result, err := tx.ExecContext(ctx, SQL, image.ProductId, image.Url)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	image.Id = int(id)
	return image
}

func (repository *ProductImageRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, imageId int) {
	SQL := "DELETE FROM product_image WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, imageId)
	helper.PanicIfError(err)
}

func (repository *ProductImageRepositoryImpl) FindByProductId(ctx context.Context, tx *sql.Tx, productId int) []domain.ProductImage {
	SQL := "SELECT id, product_id, url FROM product_image WHERE product_id = ?"
	rows, err := tx.QueryContext(ctx, SQL, productId)
	helper.PanicIfError(err)
	defer rows.Close()

	var images []domain.ProductImage
	for rows.Next() {
		image := domain.ProductImage{}
		err := rows.Scan(&image.Id, &image.ProductId, &image.Url)
		helper.PanicIfError(err)
		images = append(images, image)
	}

	return images
}
