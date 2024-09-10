package helper

import (
	"github.com/gunawan98/golang-restfull-api/model/domain"
	"github.com/gunawan98/golang-restfull-api/model/web"
)

// PRODUCT CATEGORY RESPONSE ###
func ToCategoryResponse(category domain.Category) web.CategoryResponse {
	return web.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

func ToCategoryResponses(categories []domain.Category) []web.CategoryResponse {
	var categoryResponse []web.CategoryResponse
	for _, category := range categories {
		categoryResponse = append(categoryResponse, ToCategoryResponse(category))
	}

	return categoryResponse
}

// PRODUCT RESPONSE ###
func ToProductResponse(product domain.Product) web.ProductResponse {
	return web.ProductResponse{
		Id:       product.Id,
		Name:     product.Name,
		Barcode:  product.Barcode,
		Stock:    product.Stock,
		Price:    product.Price,
		Discount: product.Discount,
	}
}

func ToProductResponses(products []domain.Product) []web.ProductResponse {
	var productResponse []web.ProductResponse
	for _, product := range products {
		productResponse = append(productResponse, ToProductResponse(product))
	}

	return productResponse
}

// USER RESPONSE ###
func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Role:     user.Role,
	}
}

func ToUserResponses(users []domain.User) []web.UserResponse {
	var userResponse []web.UserResponse
	for _, user := range users {
		userResponse = append(userResponse, ToUserResponse(user))
	}

	return userResponse
}

// CART RESPONSE ###
func ToCartResponse(cart domain.Cart) web.CartResponse {
	return web.CartResponse{
		Id:        cart.Id,
		CashierID: cart.CashierID,
		Completed: cart.Completed,
		CreatedAt: cart.CreatedAt,
	}
}

func ToCartResponses(carts []domain.Cart) []web.CartResponse {
	var cartResponse []web.CartResponse
	for _, cart := range carts {
		cartResponse = append(cartResponse, ToCartResponse(cart))
	}

	return cartResponse
}

func ToCartItemResponse(cartItem domain.CartItem) web.CartItemResponse {
	return web.CartItemResponse{
		Id:         cartItem.Id,
		CartID:     cartItem.CartID,
		ProductID:  cartItem.ProductID,
		Quantity:   cartItem.Quantity,
		UnitPrice:  cartItem.UnitPrice,
		TotalPrice: cartItem.TotalPrice,
	}
}

// PURCHASE RESPONSE ###
func ToPurchaseResponse(purchase domain.Purchase) web.PurchaseResponse {
	return web.PurchaseResponse{
		CartID:        purchase.CartID,
		CashierID:     purchase.CashierID,
		TotalAmount:   purchase.TotalAmount,
		PaymentMethod: purchase.PaymentMethod,
		CreatedAt:     purchase.CreatedAt,
	}
}
