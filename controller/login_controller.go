package controller

import (
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

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userResponse.Id,
		"username": userResponse.Username,
		"role":     userResponse.Role,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString([]byte("toko_secret_key"))
	if err != nil {
		helper.PanicIfError(err)
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   tokenString,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
