package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (repository *ProductRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	SQL := "INSERT INTO product(name, barcode, stock, price, discount) VALUES (?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, product.Name, product.Barcode, product.Stock, product.Price, product.Discount)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	product.Id = int(id)
	return product
}

func (repository *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	SQL := "UPDATE product SET name = ?, barcode = ?, stock = ?, price = ?, discount = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, product.Name, product.Barcode, product.Stock, product.Price, product.Discount, product.Id)
	helper.PanicIfError(err)

	return product
}

func (repository *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, productId int) {
	SQL := "DELETE FROM product WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, productId)
	helper.PanicIfError(err)
}

func (repository *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, productId int) (domain.Product, error) {
	SQL := "SELECT id, name, barcode, stock, price, discount FROM product WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, productId)
	helper.PanicIfError(err)
	defer rows.Close()

	product := domain.Product{}
	if rows.Next() {
		err := rows.Scan(&product.Id, &product.Name, &product.Barcode, &product.Stock, &product.Price, &product.Discount)
		helper.PanicIfError(err)
		return product, nil
	} else {
		return product, errors.New("product is not found")
	}
}

func (repository *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Product {
	SQL := "SELECT id, name, barcode, stock, price, discount FROM product"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		product := domain.Product{}
		err := rows.Scan(&product.Id, &product.Name, &product.Barcode, &product.Stock, &product.Price, &product.Discount)
		helper.PanicIfError(err)
		products = append(products, product)
	}

	return products
}

func (repository *ProductRepositoryImpl) FindByBarcode(ctx context.Context, tx *sql.Tx, barcode string) (domain.Product, error) {
	SQL := "SELECT id, name, barcode, stock, price, discount FROM product WHERE barcode = ? "
	rows, err := tx.QueryContext(ctx, SQL, barcode)
	helper.PanicIfError(err)
	defer rows.Close()

	product := domain.Product{}
	if rows.Next() {
		err := rows.Scan(&product.Id, &product.Name, &product.Barcode, &product.Stock, &product.Price, &product.Discount)
		helper.PanicIfError(err)
		return product, nil
	} else {
		return product, sql.ErrNoRows
	}
}

func (repository *ProductRepositoryImpl) UpdateStock(ctx context.Context, tx *sql.Tx, productId int, stock int) error {
	SQL := "UPDATE product SET stock = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, stock, productId)
	return err
}
