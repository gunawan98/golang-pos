package controller

import (
	"net/http"
	"strconv"

	globalctx "github.com/gunawan98/golang-restfull-api/global_ctx"
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
	userId, ok := request.Context().Value(globalctx.UserIDKey()).(float64)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	purchaseCreateRequest := web.PurchaseCreateRequest{}
	helper.ReadFromRequestBody(request, &purchaseCreateRequest)

	purchaseResponse, resCartItems := controller.PurchaseService.ConfirmPayment(request.Context(), purchaseCreateRequest, userId)

	data := map[string]interface{}{
		"purchase": purchaseResponse,
		"items":    resCartItems,
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PurchaseControllerImpl) GetFinishedPayment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, ok := request.Context().Value(globalctx.UserIDKey()).(float64)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	getCartId := params.ByName("cartId")
	cartId, err := strconv.Atoi(getCartId)
	helper.PanicIfError(err)

	purchaseResponse, resCartItems := controller.PurchaseService.GetFinishedPayment(request.Context(), userId, cartId)

	data := map[string]interface{}{
		"purchase": purchaseResponse,
		"items":    resCartItems,
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
