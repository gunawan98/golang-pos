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
	ProductRepository      repository.ProductRepository
	ProductImageRepository repository.ProductImageRepository
	DB                     *sql.DB
	Validate               *validator.Validate
}

func NewProductService(productRepository repository.ProductRepository, productImageRepository repository.ProductImageRepository, DB *sql.DB, validate *validator.Validate) ProductService {
	return &ProductServiceImpl{
		ProductRepository:      productRepository,
		ProductImageRepository: productImageRepository,
		DB:                     DB,
		Validate:               validate,
	}
}

func (service *ProductServiceImpl) Create(ctx context.Context, request web.ProductCreateRequest) web.ProductResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Checking barcode if exists
	_, err = service.ProductRepository.FindByBarcode(ctx, tx, request.Barcode)
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

	barcodeExists := service.ProductRepository.FindBarcodeOtherOwn(ctx, tx, product.Id, request.Barcode)
	if barcodeExists {
		panic(exception.NewDataAlreadyExistsError("Barcode already exists"))
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

	images := service.ProductImageRepository.FindByProductId(ctx, tx, productId)

	var imageURL []string
	for _, url := range images {
		imageURL = append(imageURL, url.Url)
	}

	// Delete the image from Cloudinary
	err = helper.DeleteMultipleFiles(imageURL)
	helper.PanicIfError(err)

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

func (service *ProductServiceImpl) FindAll(ctx context.Context, page int, perPage int) ([]web.ProductResponse, int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Fetch total count
	total := service.ProductRepository.CountAllProducts(ctx, tx)

	// Fetch paginated products
	offset := (page - 1) * perPage
	products := service.ProductRepository.FindAll(ctx, tx, perPage, offset)

	// Map products to responses
	var productResponses []web.ProductResponse
	for _, product := range products {
		images := service.ProductImageRepository.FindByProductId(ctx, tx, product.Id)
		productResponse := helper.ToProductResponseWithImages(product, images)
		productResponses = append(productResponses, productResponse)
	}

	return productResponses, total
}
