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

type CartServiceImpl struct {
	CartRepository    repository.CartRepository
	ProductRepository repository.ProductRepository
	DB                *sql.DB
	Validate          *validator.Validate
}

func NewCartService(cartRepository repository.CartRepository, productRepository repository.ProductRepository, DB *sql.DB, validate *validator.Validate) CartService {
	return &CartServiceImpl{
		CartRepository:    cartRepository,
		ProductRepository: productRepository,
		DB:                DB,
		Validate:          validate,
	}
}

func (service *CartServiceImpl) CreateNewCart(ctx context.Context, request web.CartCreateRequest) web.CartResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cart := domain.Cart{
		CashierID: request.CashierID,
		Completed: request.Completed,
	}

	cart = service.CartRepository.CreateCart(ctx, tx, cart)

	return helper.ToCartResponse(cart)
}

func (service *CartServiceImpl) AddProductToCart(ctx context.Context, cartId int, request web.CartItemCreateRequest) web.CartItemResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	_, errGetCart := service.CartRepository.GetCartById(ctx, tx, cartId)
	if errGetCart != nil {
		panic(exception.NewNotFoundError(errGetCart.Error()))
	}

	product, errGetProduct := service.ProductRepository.FindById(ctx, tx, request.ProductID)
	if errGetProduct != nil {
		panic(exception.NewNotFoundError(errGetProduct.Error()))
	}

	// Set UnitPrice and calculate TotalPrice
	unitPrice := product.Price
	totalPrice := unitPrice * request.Quantity

	cartItem := domain.CartItem{
		CartID:     cartId,
		ProductID:  request.ProductID,
		Quantity:   request.Quantity,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}

	cartItem = service.CartRepository.AddItemToCart(ctx, tx, cartItem)

	return helper.ToCartItemResponse(cartItem)
}

func (service *CartServiceImpl) GetCartDetails(ctx context.Context, cartId int) (web.CartResponse, []web.CartItemResponse) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cart, errGetCart := service.CartRepository.GetCartById(ctx, tx, cartId)
	if errGetCart != nil {
		panic(exception.NewNotFoundError(errGetCart.Error()))
	}

	cartItems := service.CartRepository.GetItemsByCartId(ctx, tx, cartId)

	var items []web.CartItemResponse
	for _, item := range cartItems {
		items = append(items, web.CartItemResponse{
			Id:         item.Id,
			CartID:     item.CartID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
		})
	}

	resCart := web.CartResponse{
		Id:        cart.Id,
		CashierID: cart.CashierID,
		Completed: cart.Completed,
		CreatedAt: cart.CreatedAt,
	}

	return resCart, items
}
