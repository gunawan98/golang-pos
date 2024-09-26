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

func (service *CartServiceImpl) AvailableCart(ctx context.Context, userId float64) []web.CartResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	listCart := service.CartRepository.FindAvailableCart(ctx, tx, int(userId))

	return helper.ToCartResponses(listCart)
}

func (service *CartServiceImpl) FinishedCart(ctx context.Context, userId float64) []web.CartResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	listCart := service.CartRepository.FindFinishedCart(ctx, tx, int(userId))

	return helper.ToCartResponses(listCart)
}

func (service *CartServiceImpl) CreateNewCart(ctx context.Context, request web.CartCreateRequest, userId float64) web.CartResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cart := domain.Cart{
		CashierID: int(userId),
		Completed: request.Completed,
	}

	cart = service.CartRepository.CreateCart(ctx, tx, cart)

	return helper.ToCartResponse(cart)
}

func (service *CartServiceImpl) AddProductToCart(ctx context.Context, userId float64, cartId int, request web.CartItemCreateRequest) web.CartItemResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	getCart, errGetCart := service.CartRepository.GetCartById(ctx, tx, cartId)
	if errGetCart != nil {
		panic(exception.NewNotFoundError(errGetCart.Error()))
	}

	if int(userId) != getCart.CashierID {
		panic(exception.NewNotFoundError("Cart ID not available"))
	}

	product, errGetProduct := service.ProductRepository.FindByBarcode(ctx, tx, request.Barcode)
	if errGetProduct != nil {
		panic(exception.NewNotFoundError(errGetProduct.Error()))
	}

	productId := product.Id
	unitPrice := product.Price
	totalPrice := unitPrice * request.Quantity

	// Cek apakah item sudah ada dalam keranjang
	cartItem, errFind := service.CartRepository.FindItemByCartAndProduct(ctx, tx, cartId, productId)
	if errFind != nil {
		if errFind == sql.ErrNoRows {
			// Jika item tidak ada, tambahkan item baru
			newCartItem := domain.CartItem{
				CartID:     cartId,
				ProductID:  productId,
				Quantity:   request.Quantity,
				UnitPrice:  unitPrice,
				TotalPrice: totalPrice,
			}
			newCartItem = service.CartRepository.AddItemToCart(ctx, tx, newCartItem)
			return helper.ToCartItemResponse(newCartItem)
		} else {
			helper.PanicIfError(errFind)
		}
	}

	// Jika item sudah ada, lakukan update
	cartItem.Quantity += request.Quantity
	cartItem.TotalPrice = cartItem.UnitPrice * cartItem.Quantity
	errUpdate := service.CartRepository.UpdateCartItem(ctx, tx, cartItem)
	helper.PanicIfError(errUpdate)

	return helper.ToCartItemResponse(cartItem)
}

func (service *CartServiceImpl) GetCartDetails(ctx context.Context, cartId int) (web.CartResponse, []web.CartItemWithProductResponse) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cart, errGetCart := service.CartRepository.GetCartById(ctx, tx, cartId)
	if errGetCart != nil {
		panic(exception.NewNotFoundError(errGetCart.Error()))
	}

	cartItems := service.CartRepository.GetItemsWithProductByCartId(ctx, tx, cartId)

	var items []web.CartItemWithProductResponse
	for _, item := range cartItems {
		items = append(items, web.CartItemWithProductResponse{
			Id:          item.Id,
			CartID:      item.CartID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
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
