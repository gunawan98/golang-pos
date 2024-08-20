package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gunawan98/golang-restfull-api/app"
	"github.com/gunawan98/golang-restfull-api/controller"
	"github.com/gunawan98/golang-restfull-api/exception"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/middleware"
	"github.com/gunawan98/golang-restfull-api/repository"
	"github.com/gunawan98/golang-restfull-api/service"
	"github.com/julienschmidt/httprouter"
)

func main() {

	db := app.NewDB()
	validate := validator.New()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate)
	userController := controller.NewUserController(userService)

	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(productRepository, db, validate)
	productController := controller.NewProductController(productService)

	cartRepository := repository.NewCartRepository()
	cartService := service.NewCartService(cartRepository, productRepository, db, validate)
	cartController := controller.NewCartController(cartService)

	purchaseRepository := repository.NewPurchaseRepository()
	purchaseService := service.NewPurchaseService(purchaseRepository, cartRepository, productRepository, db, validate)
	purchaseController := controller.NewPurchaseController(purchaseService)

	loginController := controller.NewLoginController(userService)

	router := httprouter.New()

	router.POST("/api/login", loginController.Login)
	router.POST("/api/refresh", loginController.Refresh)

	router.GET("/api/categories", categoryController.FindAll)
	router.GET("/api/categories/:categoryId", categoryController.FindById)
	router.POST("/api/categories", categoryController.Create)
	router.PUT("/api/categories/:categoryId", categoryController.Update)
	router.DELETE("/api/categories/:categoryId", categoryController.Delete)

	router.GET("/api/product", productController.FindAll)
	router.GET("/api/product/:productId", productController.FindById)
	router.POST("/api/product", productController.Create)
	router.PUT("/api/product/:productId", productController.Update)
	router.DELETE("/api/product/:productId", productController.Delete)

	router.GET("/api/user", userController.FindAll)
	router.GET("/api/user/:userId", userController.FindById)
	router.POST("/api/user", userController.Create)
	router.PUT("/api/user/:userId", userController.Update)
	router.DELETE("/api/user/:userId", userController.Delete)

	router.POST("/api/cart", cartController.CreateCart)
	router.POST("/api/cart-item/:cartId", cartController.AddItem)
	router.GET("/api/cart-item/:cartId", cartController.GetCartDetails)

	router.POST("/api/purchase", purchaseController.ConfirmPayment)

	router.PanicHandler = exception.ErrorHandler

	// Protect routes with the middleware
	protectedRouter := middleware.NewAuthMiddleware(router)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: protectedRouter,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
