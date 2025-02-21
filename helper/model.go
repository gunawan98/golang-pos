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

func ToProductResponseWithImages(product domain.Product, images []domain.ProductImage) web.ProductResponse {
	imageUrls := make([]string, len(images))
	for i, image := range images {
		imageUrls[i] = image.Url
	}

	return web.ProductResponse{
		Id:       product.Id,
		Name:     product.Name,
		Barcode:  product.Barcode,
		Stock:    product.Stock,
		Price:    product.Price,
		Discount: product.Discount,
		Images:   imageUrls,
	}
}

func ToProductResponses(products []domain.Product, imageMap map[int][]domain.ProductImage) []web.ProductResponse {
	var productResponses []web.ProductResponse
	for _, product := range products {
		images := imageMap[product.Id]
		productResponse := ToProductResponseWithImages(product, images)
		productResponses = append(productResponses, productResponse)
	}

	return productResponses
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
		Paid:          purchase.Paid,
		CashBack:      purchase.CashBack,
		PaymentMethod: purchase.PaymentMethod,
		CreatedAt:     purchase.CreatedAt,
	}
}

// PRODUCT IMAGE RESPONSE ###
func ToImageResponse(image domain.ProductImage) web.ProductImageResponse {
	return web.ProductImageResponse{
		Id:        image.Id,
		ProductId: image.ProductId,
		Url:       image.Url,
	}
}

func ToImageResponses(images []domain.ProductImage) []web.ProductImageResponse {
	var imageResponses []web.ProductImageResponse
	for _, image := range images {
		imageResponses = append(imageResponses, ToImageResponse(image))
	}

	return imageResponses
}
