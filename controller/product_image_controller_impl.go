package controller

import (
	"net/http"
	"strconv"

	"github.com/gunawan98/golang-restfull-api/helper"
	"github.com/gunawan98/golang-restfull-api/model/web"
	"github.com/gunawan98/golang-restfull-api/service"
	"github.com/julienschmidt/httprouter"
)

type ProductImageControllerImpl struct {
	ProductImageService service.ProductImageService
}

func NewProductImageController(imageService service.ProductImageService) ProductImageController {
	return &ProductImageControllerImpl{
		ProductImageService: imageService,
	}
}

func (controller *ProductImageControllerImpl) AddImage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := request.ParseMultipartForm(10 << 20) // max 10MB
	helper.PanicIfError(err)

	file, _, err := request.FormFile("file")
	helper.PanicIfError(err)
	defer file.Close()

	productId, err := strconv.Atoi(params.ByName("productId"))
	helper.PanicIfError(err)

	filePath := helper.SaveFile(file)

	imageCreateRequest := web.ProductImageCreateRequest{
		ProductId: productId,
		Url:       filePath,
	}

	imageResponse := controller.ProductImageService.AddImage(request.Context(), imageCreateRequest)
	webResponse := web.WebResponse{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    imageResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProductImageControllerImpl) DeleteImage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	imageId := params.ByName("imageId")
	id, err := strconv.Atoi(imageId)
	helper.PanicIfError(err)

	// Retrieve the image URL from the database
	imageResponse := controller.ProductImageService.FindById(request.Context(), id)
	imageURL := imageResponse.Url

	// Delete the image from Cloudinary
	err = helper.DeleteFile(imageURL)
	helper.PanicIfError(err)

	// Delete the image record from the database
	controller.ProductImageService.DeleteImage(request.Context(), id)
	webResponse := web.WebResponse{
		Success: true,
		Code:    200,
		Message: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProductImageControllerImpl) FindByProductId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId := params.ByName("productId")
	id, err := strconv.Atoi(productId)
	helper.PanicIfError(err)

	imageResponses := controller.ProductImageService.FindByProductId(request.Context(), id)
	webResponse := web.WebResponse{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    imageResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
