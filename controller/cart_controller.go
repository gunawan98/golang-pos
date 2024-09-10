package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CartController interface {
	AvailableCart(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateCart(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	AddItem(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetCartDetails(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
