package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
	// Bypass authentication for login endpoint
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
