package app

import (
	"github.com/gunawan98/golang-restfull-api/controller"
	"github.com/gunawan98/golang-restfull-api/exception"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(loginController controller.LoginController, categoryController controller.CategoryController, productController controller.ProductController, userController controller.UserController, cartController controller.CartController, purchaseController controller.PurchaseController) *httprouter.Router {
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

	return router
}