package controller

import (
	"net/http"

	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/web"
	"github.com/gunawan98/golang-restfull-api/service"
	"github.com/julienschmidt/httprouter"
)

type PurchaseControllerImpl struct {
	PurchaseService service.PurchaseService
}

func NewPurchaseController(purchaseService service.PurchaseService) PurchaseController {
	return &PurchaseControllerImpl{
		PurchaseService: purchaseService,
	}
}

func (controller *PurchaseControllerImpl) ConfirmPayment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	purchaseCreateRequest := web.PurchaseCreateRequest{}
	helper.ReadFromRequestBody(request, &purchaseCreateRequest)

	purchaseResponse := controller.PurchaseService.ConfirmPayment(request.Context(), purchaseCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   purchaseResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
