package controller

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/web"
	"github.com/julienschmidt/httprouter"
)

type LoginController struct {
}

func NewLoginController() *LoginController {
	return &LoginController{}
}

func (controller *LoginController) Login(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var credentials web.LoginRequest
	helper.ReadFromRequestBody(request, &credentials)

	if credentials.Username == "user" && credentials.Password == "password" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": credentials.Username,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		})

		tokenString, err := token.SignedString([]byte("your_secret_key"))
		if err != nil {
			helper.PanicIfError(err)
			return
		}

		response := web.WebResponse{
			Code:   http.StatusOK,
			Status: "OK",
			Data:   tokenString,
		}

		helper.WriteToResponseBody(writer, response)
	} else {
		response := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}

		helper.WriteToResponseBody(writer, response)
	}
}
