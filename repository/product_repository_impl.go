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
	SQL := "insert into product(name, barcode, stock, price, discount) values (?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, product.Name, product.Barcode, product.Stock, product.Price, product.Discount)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	product.Id = int(id)
	return product
}

func (repository *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	SQL := "update product set name = ?, barcode = ?, stock = ?, price = ?, discount = ? where id = ?"
	_, err := tx.ExecContext(ctx, SQL, product.Name, product.Barcode, product.Stock, product.Price, product.Discount, product.Id)
	helper.PanicIfError(err)

	return product
}

func (repository *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, productId int) {
	SQL := "delete from product where id = ?"
	_, err := tx.ExecContext(ctx, SQL, productId)
	helper.PanicIfError(err)
}

func (repository *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, productId int) (domain.Product, error) {
	SQL := "select id, name, barcode, stock, price, discount from product where id = ?"
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
	SQL := "select id, name, barcode, stock, price, discount from product"
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

func (repository *ProductRepositoryImpl) FindByBarcode(ctx context.Context, tx *sql.Tx, id int, barcode string) (domain.Product, error) {
	SQL := "select id, name, barcode, stock, price, discount from product where barcode = ? and id != ?"
	rows, err := tx.QueryContext(ctx, SQL, barcode, id)
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
