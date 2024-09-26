package service

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gunawan98/golang-restfull-api/exception"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/domain"
	"github.com/gunawan98/golang-restfull-api/model/web"
	"github.com/gunawan98/golang-restfull-api/repository"
)

type PurchaseServiceImpl struct {
	PurchaseRepository repository.PurchaseRepository
	CartRepository     repository.CartRepository
	ProductRepository  repository.ProductRepository
	DB                 *sql.DB
	Validate           *validator.Validate
}

func NewPurchaseService(purchaseRepository repository.PurchaseRepository, cartRepository repository.CartRepository, productRepository repository.ProductRepository, DB *sql.DB, validate *validator.Validate) PurchaseService {
	return &PurchaseServiceImpl{
		PurchaseRepository: purchaseRepository,
		CartRepository:     cartRepository,
		ProductRepository:  productRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

func (service *PurchaseServiceImpl) ConfirmPayment(ctx context.Context, request web.PurchaseCreateRequest, userId float64) (web.PurchaseResponse, []web.CartItemWithProductResponse) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Get cart by ID to check if it is completed
	cart, err := service.CartRepository.GetCartById(ctx, tx, request.CartID)
	helper.PanicIfError(err)

	// Validate if the cart is already completed
	if cart.Completed {
		panic(exception.NewBadRequestError("Transaction has been completed previously"))
	}

	// Retrieve cart items
	cartItems := service.CartRepository.GetItemsWithProductByCartId(ctx, tx, request.CartID)

	// Update stock for each product
	var resCartItems []web.CartItemWithProductResponse
	for _, item := range cartItems {
		resCartItems = append(resCartItems, web.CartItemWithProductResponse{
			Id:          item.Id,
			CartID:      item.CartID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		})

		product, err := service.ProductRepository.FindById(ctx, tx, item.ProductID)
		helper.PanicIfError(err)

		newStock := product.Stock - item.Quantity
		if newStock < 0 {
			panic(exception.NewBadRequestError("Insufficient stock for product id " + strconv.Itoa(item.ProductID)))
		}
		err = service.ProductRepository.UpdateStock(ctx, tx, item.ProductID, newStock)
		helper.PanicIfError(err)
	}

	// Update cart status
	err = service.CartRepository.UpdateCartStatus(ctx, tx, request.CartID)
	helper.PanicIfError(err)

	// Calculate total amount
	totalAmount := 0
	for _, item := range cartItems {
		totalAmount += item.TotalPrice
	}
	if totalAmount > request.Paid {
		panic(exception.NewBadRequestError("Money is not enough"))
	}

	// Calculate cashback
	cashBack := 0
	if request.Paid > totalAmount {
		cashBack = request.Paid - totalAmount
	}

	// Create purchase record
	purchase := domain.Purchase{
		CartID:        request.CartID,
		CashierID:     int(userId),
		TotalAmount:   totalAmount,
		Paid:          request.Paid,
		CashBack:      cashBack,
		PaymentMethod: request.PaymentMethod,
		CreatedAt:     time.Now(),
	}
	purchase = service.PurchaseRepository.AddPurchase(ctx, tx, purchase)

	return helper.ToPurchaseResponse(purchase), resCartItems
}

func (service *PurchaseServiceImpl) GetFinishedPayment(ctx context.Context, userId float64, cartId int) (web.PurchaseResponse, []web.CartItemWithProductResponse) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Get cart by ID to check if it is completed
	cart, err := service.CartRepository.GetCartById(ctx, tx, cartId)
	helper.PanicIfError(err)

	// Validate if the cart is already completed
	if !cart.Completed {
		panic(exception.NewBadRequestError("Cart has not been paid"))
	}

	// Retrieve cart items
	cartItems := service.CartRepository.GetItemsWithProductByCartId(ctx, tx, cartId)

	// Update stock for each product
	var resCartItems []web.CartItemWithProductResponse
	for _, item := range cartItems {
		resCartItems = append(resCartItems, web.CartItemWithProductResponse{
			Id:          item.Id,
			CartID:      item.CartID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		})
	}

	purchase, errGetPurchase := service.PurchaseRepository.GetPurchaseByCartId(ctx, tx, int(userId), cartId)
	if errGetPurchase != nil {
		panic(exception.NewNotFoundError(errGetPurchase.Error()))
	}

	return helper.ToPurchaseResponse(purchase), resCartItems
}
