package app

import (
	"fmt"
	"net/http"

	"github.com/gunawan98/golang-restfull-api/controller"
	"github.com/gunawan98/golang-restfull-api/exception"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome To test production!\n")
}

func NewRouter(loginController controller.LoginController, categoryController controller.CategoryController, productController controller.ProductController, userController controller.UserController, cartController controller.CartController, purchaseController controller.PurchaseController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/", Index)
	router.POST("/api/login", loginController.Login)
	router.POST("/api/refresh", loginController.Refresh)

	router.GET("/api/category", categoryController.FindAll)
	router.GET("/api/category/:categoryId", categoryController.FindById)
	router.POST("/api/category", categoryController.Create)
	router.PUT("/api/category/:categoryId", categoryController.Update)
	router.DELETE("/api/category/:categoryId", categoryController.Delete)

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

	router.GET("/api/cart", cartController.AvailableCart)
	router.GET("/api/cart/finished", cartController.FinishedCart)
	router.POST("/api/cart", cartController.CreateCart)
	router.POST("/api/cart-item/:cartId", cartController.AddItem)
	router.GET("/api/cart-item/:cartId", cartController.GetCartDetails)

	router.GET("/api/purchase/:cartId", purchaseController.GetFinishedPayment)
	router.POST("/api/purchase", purchaseController.ConfirmPayment)

	router.PanicHandler = exception.ErrorHandler

	return router
}
