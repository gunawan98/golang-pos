package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gunawan98/golang-restfull-api/app"
	"github.com/gunawan98/golang-restfull-api/config"
	"github.com/gunawan98/golang-restfull-api/controller"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/middleware"
	"github.com/gunawan98/golang-restfull-api/middleware/log"
	"github.com/gunawan98/golang-restfull-api/repository"
	"github.com/gunawan98/golang-restfull-api/service"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config.LoadEnvVars()

	log.LoadLogger()
	db := config.MySQLConnect()

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

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	}).Handler(protectedRouter)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}

	log.Logger.Info("Server is listening on http://localhost: " + port)
	// fmt.Printf("Server is listening on http://localhost:%s\n", port)

	server := http.Server{
		Addr: ":" + port,
		// Handler: corsHandler,
		Handler: WrapHandler(corsHandler),
		// ErrorLog:WrapHandler,
	}

	err := server.ListenAndServe()

	helper.PanicIfError(err)
}

type CustomResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // Capture the response body in the buffer
	return w.ResponseWriter.Write(b)
}

func WrapHandler(f http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Clone the request body to avoid losing it after reading
		var reqBodyClone []byte
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				log.Logger.Error("Failed to read request body: " + err.Error())
			}
			reqBodyClone = bodyBytes

			// Restore the request body so the next handler can read it
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Use CustomResponseWriter to capture the response body
		crw := &CustomResponseWriter{
			ResponseWriter: res,
			body:           bytes.NewBuffer(nil),
		}

		// Call the next handler
		f.ServeHTTP(crw, req)

		// Log the request and response
		responseTime := time.Now()
		requestTime := time.Now()

		// set request response
		fields := []zapcore.Field{
			zap.String("unique_id", uuid.New().String()),
			zap.String("request", string(reqBodyClone)),
			zap.String("response", crw.body.String()),
		}
		if req != nil {
			fields = append(fields, zap.String("request_time", requestTime.String()))
		}

		if res != nil {
			fields = append(fields, zap.String("response_time", responseTime.String()))
			processingTime := time.Since(requestTime)
			fields = append(fields, zap.Int("processing_time_nano_second", int(processingTime.Nanoseconds())))
		}
		log.Logger.Info("log global request and response", fields...)
	}
}
