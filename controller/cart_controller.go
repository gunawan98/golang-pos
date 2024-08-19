package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CartController interface {
	CreateCart(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	AddItem(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// GetCartDetails(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
