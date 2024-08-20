package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PurchaseController interface {
	ConfirmPayment(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
