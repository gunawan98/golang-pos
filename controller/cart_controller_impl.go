package controller

import (
	"net/http"
	"strconv"

	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/web"
	"github.com/gunawan98/golang-restfull-api/service"
	"github.com/julienschmidt/httprouter"
)

type CartControllerImpl struct {
	CartService service.CartService
}

func NewCartController(cartService service.CartService) CartController {
	return &CartControllerImpl{CartService: cartService}
}

func (controller *CartControllerImpl) CreateCart(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	cartCreateRequest := web.CartCreateRequest{}
	helper.ReadFromRequestBody(request, &cartCreateRequest)

	cartResponse := controller.CartService.CreateNewCart(request.Context(), cartCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   cartResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CartControllerImpl) AddItem(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	cartItemCreateRequest := web.CartItemCreateRequest{}
	helper.ReadFromRequestBody(request, &cartItemCreateRequest)

	getCartId := params.ByName("cartId")
	cartId, err := strconv.Atoi(getCartId)
	helper.PanicIfError(err)

	cartResponse := controller.CartService.AddProductToCart(request.Context(), cartId, cartItemCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   cartResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CartControllerImpl) GetCartDetails(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	getCartId := params.ByName("cartId")
	cartId, err := strconv.Atoi(getCartId)
	helper.PanicIfError(err)

	cart, items := controller.CartService.GetCartDetails(request.Context(), cartId)

	data := map[string]interface{}{
		"Cart":  cart,
		"Items": items,
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.WriteToResponseBody(writer, webResponse)

}
