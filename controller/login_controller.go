package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/web"
	"github.com/gunawan98/golang-restfull-api/service"
	"github.com/julienschmidt/httprouter"
)

type LoginController struct {
	UserService service.UserService
}

func NewLoginController(userService service.UserService) *LoginController {
	return &LoginController{
		UserService: userService,
	}
}

func (controller *LoginController) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	loginRequest := web.LoginRequest{}
	helper.ReadFromRequestBody(request, &loginRequest)

	userResponse := controller.UserService.Authenticate(request.Context(), loginRequest)
	// if err != nil {
	// 		webResponse := web.WebResponse{
	// 				Code:   http.StatusUnauthorized,
	// 				Status: "Unauthorized",
	// 		}
	// 		helper.WriteToResponseBody(writer, webResponse)
	// 		return
	// }

	// Define expiration times
	accessExpirationTime := time.Now().Add(60 * time.Minute)
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days

	// Generate Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userResponse.Id,
		"username": userResponse.Username,
		"exp":      accessExpirationTime.Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte("your_secret_key"))
	if err != nil {
		helper.PanicIfError(err)
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userResponse.Id,
		"username": userResponse.Username,
		"exp":      refreshExpirationTime.Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte("your_refresh_secret_key"))
	if err != nil {
		helper.PanicIfError(err)
	}

	// Create the response
	tokenResponse := web.TokenResponse{
		AccessToken:       accessTokenString,
		AccessValidUntil:  accessExpirationTime.Format(time.RFC3339),
		RefreshToken:      refreshTokenString,
		RefreshValidUntil: refreshExpirationTime.Format(time.RFC3339),
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   tokenResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *LoginController) Refresh(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	refreshRequest := struct {
		RefreshToken string `json:"refresh"`
	}{}
	helper.ReadFromRequestBody(request, &refreshRequest)

	fmt.Println("Received refresh token:", refreshRequest.RefreshToken)

	// Parse the token
	token, err := jwt.Parse(refreshRequest.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_refresh_secret_key"), nil
	})

	if err != nil {
		fmt.Println("Token parsing error:", err)
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: fmt.Sprintf("Invalid refresh token: %v", err),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Println("Invalid token or claims extraction failed")
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "Invalid token claims",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	userId := claims["userId"]
	username := claims["username"]

	// Generate new access token
	accessExpirationTime := time.Now().Add(60 * time.Minute)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      accessExpirationTime.Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte("your_secret_key"))
	if err != nil {
		fmt.Println("Error signing new access token:", err)
		helper.PanicIfError(err)
	}

	// Generate new refresh token
	refreshExpirationTime := time.Now().Add(24 * time.Hour)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      refreshExpirationTime.Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte("your_refresh_secret_key"))
	if err != nil {
		fmt.Println("Error signing new refresh token:", err)
		helper.PanicIfError(err)
	}

	// Create the response
	tokenResponse := map[string]string{
		"access":              accessTokenString,
		"access_valid_until":  accessExpirationTime.Format(time.RFC3339),
		"refresh":             refreshTokenString,
		"refresh_valid_until": refreshExpirationTime.Format(time.RFC3339),
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   tokenResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
