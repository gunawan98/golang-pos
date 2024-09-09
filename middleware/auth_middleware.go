package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	globalctx "github.com/gunawan98/golang-restfull-api/global_ctx"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/web"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Bypass authentication for login and refresh endpoints
	if request.URL.Path == "/api/login" || request.URL.Path == "/api/refresh" {
		middleware.Handler.ServeHTTP(writer, request)
		return
	}

	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(writer, http.StatusUnauthorized, "Missing authorization header")
		return
	}

	tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("your_secret_key"), nil
	})

	if err != nil || !token.Valid {
		respondWithError(writer, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respondWithError(writer, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	userIdFloat, ok := claims["userId"].(float64)
	if !ok {
		respondWithError(writer, http.StatusUnauthorized, "Invalid token userId")
		return
	}

	role, ok := claims["role"].(string)
	if !ok {
		respondWithError(writer, http.StatusUnauthorized, "Invalid token role")
		return
	}

	requiresAdmin := map[string]bool{
		"/api/user":     true,
		"/api/category": true,
		"/api/product":  true,
	}

	if requiresAdmin[request.URL.Path] {
		if role != "admin" {
			respondWithError(writer, http.StatusForbidden, "Access denied: insufficient role")
			return
		}
	}

	requiresAdminOrCashier := map[string]bool{
		"/api/cart":      true,
		"/api/cart-item": true,
		"/api/purchase":  true,
	}

	if requiresAdminOrCashier[request.URL.Path] {
		if role != "admin" && role != "cashier" {
			respondWithError(writer, http.StatusForbidden, "Access denied: insufficient role")
			return
		}
	}

	ctx := context.WithValue(request.Context(), globalctx.UserIDKey(), userIdFloat)
	request = request.WithContext(ctx)

	// Allow the request to proceed
	middleware.Handler.ServeHTTP(writer, request)
}

func respondWithError(writer http.ResponseWriter, code int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	webResponse := web.WebResponse{
		Code:   code,
		Status: message,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
