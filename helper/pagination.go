package helper

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gunawan98/golang-restfull-api/model/web"
)

// ParsePagination extracts pagination parameters from the request
func ParsePagination(request *http.Request) (int, int) {
	query := request.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	perPage, _ := strconv.Atoi(query.Get("per_page"))

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	return page, perPage
}

// CreatePaginatedResponse creates a response with metadata
func CreatePaginatedResponse(data interface{}, total int, page int, perPage int) web.WebResponse {
	// Safely determine the length of the data
	var count int
	switch v := data.(type) {
	case []web.UserResponse:
		count = len(v)
	case []web.ProductResponse:
		count = len(v)
	case []web.CartResponse:
		count = len(v)
	case []interface{}:
		count = len(v)
	case []string:
		count = len(v)
	case []int:
		count = len(v)
	case []float64:
		count = len(v)
	default:
		panic(fmt.Sprintf("Unsupported data type for pagination: %T", data))
	}

	// Create metadata
	metadata := map[string]interface{}{
		"total":        total,
		"count":        count,
		"page":         page,
		"per_page":     perPage,
		"has_next":     page*perPage < total,
		"has_previous": page > 1,
	}

	// Return the response
	return web.WebResponse{
		Success:  true,
		Code:     200,
		Message:  "Resources retrieved successfully",
		Data:     data,
		Metadata: metadata,
	}
}
