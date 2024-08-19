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

// func (c *CartControllerImpl) GetCartDetails(w http.ResponseWriter, r *http.Request) {
// 	cartId, err := strconv.Atoi(r.URL.Query().Get("cart_id"))
// 	if err != nil {
// 		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
// 		return
// 	}

// 	cart, items, err := c.cartService.GetCartDetails(cartId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response := struct {
// 		Cart  *domain.Cart      `json:"cart"`
// 		Items []domain.CartItem `json:"items"`
// 	}{
// 		Cart:  cart,
// 		Items: items,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
