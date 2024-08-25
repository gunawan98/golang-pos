package main

import (
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gunawan98/golang-restfull-api/app"
	"github.com/gunawan98/golang-restfull-api/config"
	"github.com/gunawan98/golang-restfull-api/controller"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/middleware"
	"github.com/gunawan98/golang-restfull-api/repository"
	"github.com/gunawan98/golang-restfull-api/service"
)

func main() {
	config.LoadEnvVars()
	db := config.MySQLConnect()

	// db := app.NewDB()
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

	router := app.NewRouter(loginController, categoryController, productController, userController, cartController, purchaseController)

	// Protect routes with the middleware
	protectedRouter := middleware.NewAuthMiddleware(router)

	// Get the port from environment variables (Koyeb will set this)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}
	// fmt.Print(port)
	server := http.Server{
		Addr:    ":" + port,
		Handler: protectedRouter,
	}

	// server := http.Server{
	// 	Addr:    "localhost:3000",
	// 	Handler: protectedRouter,
	// }

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
