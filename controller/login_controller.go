package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LoginController interface {
	Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Refresh(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
