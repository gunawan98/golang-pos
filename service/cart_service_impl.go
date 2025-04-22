package service

import (
	"context"
	"database/sql"
	"fmt"

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

func (service *CartServiceImpl) FinishedCart(ctx context.Context, userId float64, page int, perPage int) ([]web.CartResponse, int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	total := service.CartRepository.CountAllCartIsFinished(ctx, tx)
	offset := (page - 1) * perPage
	listCart := service.CartRepository.FindFinishedCart(ctx, tx, int(userId), perPage, offset)

	return helper.ToCartResponses(listCart), total
}

func (service *CartServiceImpl) CreateNewCart(ctx context.Context, userId float64) web.CartResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cart := domain.Cart{
		CashierID: int(userId),
		Completed: false,
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

func (service *CartServiceImpl) UpdateProductInCart(ctx context.Context, userId float64, cartId int, request web.CartItemUpdateRequest) web.CartItemResponse {
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

	// Cek apakah item sudah ada dalam keranjang
	cartItem, errFind := service.CartRepository.FindItemByCartAndProduct(ctx, tx, cartId, request.ProductID)
	if errFind == nil {
		cartItem.Quantity = request.Quantity
		cartItem.TotalPrice = cartItem.UnitPrice * cartItem.Quantity
		errUpdate := service.CartRepository.UpdateCartItem(ctx, tx, cartItem)
		helper.PanicIfError(errUpdate)
		fmt.Println("Halooooo: ", cartItem)
	} else {
		helper.PanicIfError(errFind)
	}

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
	var totalPurchase int
	for _, item := range cartItems {
		items = append(items, web.CartItemWithProductResponse{
			Id:           item.Id,
			CartID:       item.CartID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.TotalPrice,
		})
		totalPurchase += item.TotalPrice
	}

	resCart := web.CartResponse{
		Id:            cart.Id,
		CashierID:     cart.CashierID,
		Completed:     cart.Completed,
		CreatedAt:     cart.CreatedAt,
		TotalPurchase: totalPurchase,
	}

	return resCart, items
}

func (service *CartServiceImpl) DeleteCart(ctx context.Context, userId float64, cartId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cart, err := service.CartRepository.GetCartById(ctx, tx, cartId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if int(userId) != cart.CashierID {
		panic(exception.NewBadRequestError("Access denied - forbidden"))
	}

	if cart.Completed {
		panic(exception.NewBadRequestError("Cannot delete a completed cart"))
	}

	service.CartRepository.DeleteCart(ctx, tx, cart.Id)
}

func (service *CartServiceImpl) DeleteCartItem(ctx context.Context, userId float64, cartItemId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// cart, err := service.CartRepository.GetCartById(ctx, tx, cartItemId)
	// if err != nil {
	// 	panic(exception.NewNotFoundError(err.Error()))
	// }

	// if int(userId) != cart.CashierID {
	// 	panic(exception.NewBadRequestError("Access denied - forbidden"))
	// }

	service.CartRepository.DeleteCartItem(ctx, tx, cartItemId)
}
