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

func (service *PurchaseServiceImpl) ConfirmPayment(ctx context.Context, request web.PurchaseCreateRequest) web.PurchaseResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Retrieve cart items
	cartItems := service.CartRepository.GetItemsByCartId(ctx, tx, request.CartID)

	// Update stock for each product
	for _, item := range cartItems {
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

	// Create purchase record
	purchase := domain.Purchase{
		CartID:        request.CartID,
		CashierID:     request.CashierID,
		TotalAmount:   totalAmount,
		PaymentMethod: request.PaymentMethod,
		CreatedAt:     time.Now(),
	}
	purchase = service.PurchaseRepository.AddPurchase(ctx, tx, purchase)

	return helper.ToPurchaseResponse(purchase)
}
