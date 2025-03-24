package exception

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/web"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	if notFoundError(writer, request, err) {
		return
	}

	if validationError(writer, request, err) {
		return
	}

	if badRequestError(writer, request, err) {
		return
	}

	if dataAlreadyExistsError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}

func validationError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		// Customize error message for PaymentMethod field
		var customErrors []string
		for _, e := range exception {
			if e.Field() == "Role" {
				customErrors = append(customErrors, "Invalid role user. Accepted values are: user, cashier, admin.")
			} else if e.Field() == "PaymentMethod" {
				customErrors = append(customErrors, "Invalid payment method. Accepted values are: cash, credit-card, ewallet, other.")
			} else {
				customErrors = append(customErrors, e.Error()) // default error message
			}
		}

		webResponse := web.WebResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "Bad request. Please check your request body.",
			Data:    customErrors,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func dataAlreadyExistsError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(DataAlreadyExistsError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)

		webResponse := web.WebResponse{
			Success: false,
			Code:    http.StatusConflict,
			Message: "Conflict data",
			Data:    exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func badRequestError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(BadRequestError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "Bad request. Please check your request body.",
			Data:    exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func notFoundError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Success: false,
			Code:    http.StatusNotFound,
			Message: "Resource not found.",
			Data:    exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, _ *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	webResponse := web.WebResponse{
		Success: false,
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		Data:    err,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
