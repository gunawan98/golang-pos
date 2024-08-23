package test

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gunawan98/golang-restfull-api/app"
	"github.com/gunawan98/golang-restfull-api/controller"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/middleware"
	"github.com/gunawan98/golang-restfull-api/repository"
	"github.com/gunawan98/golang-restfull-api/service"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:gunawan98@tcp(localhost:3306)/retail_pos_test?parseTime=true")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
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

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
}
