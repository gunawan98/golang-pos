package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ProductImageController interface {
	AddImage(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteImage(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByProductId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
