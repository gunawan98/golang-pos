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

type ProductImageServiceImpl struct {
	ProductImageRepository repository.ProductImageRepository
	DB                     *sql.DB
	Validate               *validator.Validate
}

func NewProductImageService(imageRepository repository.ProductImageRepository, DB *sql.DB, validate *validator.Validate) ProductImageService {
	return &ProductImageServiceImpl{
		ProductImageRepository: imageRepository,
		DB:                     DB,
		Validate:               validate,
	}
}

func (service *ProductImageServiceImpl) AddImage(ctx context.Context, request web.ProductImageCreateRequest) web.ProductImageResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	image := domain.ProductImage{
		ProductId: request.ProductId,
		Url:       request.Url,
	}

	image = service.ProductImageRepository.Save(ctx, tx, image)

	return helper.ToImageResponse(image)
}

func (service *ProductImageServiceImpl) DeleteImage(ctx context.Context, imageId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Check if the image exists
	_, err = service.ProductImageRepository.FindById(ctx, tx, imageId)
	if err != nil {
		if err == sql.ErrNoRows {
			panic(exception.NewNotFoundError("Image not found"))
		} else {
			panic(err)
		}
	}

	service.ProductImageRepository.Delete(ctx, tx, imageId)
}

func (service *ProductImageServiceImpl) FindByProductId(ctx context.Context, productId int) []web.ProductImageResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	images := service.ProductImageRepository.FindByProductId(ctx, tx, productId)

	// Ensure the response is an empty slice if no images are found
	if images == nil {
		return []web.ProductImageResponse{} // Initialize as an empty slice
	}

	return helper.ToImageResponses(images)
}

func (service *ProductImageServiceImpl) FindById(ctx context.Context, imageId int) web.ProductImageResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	image, err := service.ProductImageRepository.FindById(ctx, tx, imageId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToImageResponse(image)
}
