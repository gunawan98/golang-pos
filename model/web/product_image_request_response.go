package web

type ProductImageCreateRequest struct {
	ProductId int    `validate:"required" json:"product_id"`
	Url       string `validate:"required" json:"url"`
}

type ProductImageResponse struct {
	Id        int    `json:"id"`
	ProductId int    `json:"product_id"`
	Url       string `json:"url"`
}
