package service

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/gunawan98/golang-restfull-api/exception"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/domain"
	"github.com/gunawan98/golang-restfull-api/model/web"
	"github.com/gunawan98/golang-restfull-api/repository"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
	DB                *sql.DB
	Validate          *validator.Validate
}

func NewProductService(productRepository repository.ProductRepository, DB *sql.DB, validate *validator.Validate) ProductService {
	return &ProductServiceImpl{
		ProductRepository: productRepository,
		DB:                DB,
		Validate:          validate,
	}
}

func (service *ProductServiceImpl) Create(ctx context.Context, request web.ProductCreateRequest) web.ProductResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Checking barcode if exists
	_, err = service.ProductRepository.FindByBarcode(ctx, tx, 0, request.Barcode)
	if err == nil {
		panic(exception.NewDataAlreadyExistsError("Barcode already exists"))
	} else if err != sql.ErrNoRows {
		helper.PanicIfError(err)
	}

	product := domain.Product{
		Name:     request.Name,
		Barcode:  request.Barcode,
		Stock:    request.Stock,
		Price:    request.Price,
		Discount: request.Discount,
	}

	product = service.ProductRepository.Save(ctx, tx, product)

	return helper.ToProductResponse(product)
}

func (service *ProductServiceImpl) Update(ctx context.Context, request web.ProductUpdateRequest) web.ProductResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = service.ProductRepository.FindByBarcode(ctx, tx, request.Id, request.Barcode)
	if err == nil {
		panic(exception.NewDataAlreadyExistsError("Barcode already exists"))
	} else if err != sql.ErrNoRows {
		helper.PanicIfError(err)
	}

	product.Name = request.Name
	product.Barcode = request.Barcode
	product.Stock = request.Stock
	product.Price = request.Price
	product.Discount = request.Discount

	product = service.ProductRepository.Update(ctx, tx, product)

	return helper.ToProductResponse(product)
}

func (service *ProductServiceImpl) Delete(ctx context.Context, productId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindById(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.ProductRepository.Delete(ctx, tx, product.Id)
}

func (service *ProductServiceImpl) FindById(ctx context.Context, productId int) web.ProductResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindById(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToProductResponse(product)
}

func (service *ProductServiceImpl) FindAll(ctx context.Context) []web.ProductResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	products := service.ProductRepository.FindAll(ctx, tx)

	return helper.ToProductResponses(products)
}
