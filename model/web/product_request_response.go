package web

type ProductCreateRequest struct {
	Name     string `validate:"required,max=250,min=1" json:"name"`
	Barcode  string `validate:"required,max=100,min=1" json:"barcode"`
	Stock    int    `validate:"required" json:"stock"`
	Price    int    `validate:"required" json:"price"`
	Discount *int   `json:"discount"`
}

type ProductUpdateRequest struct {
	Id       int    `validate:"required"`
	Name     string `validate:"required,max=250,min=1" json:"name"`
	Barcode  string `validate:"required,max=100,min=1" json:"barcode"`
	Stock    int    `validate:"required" json:"stock"`
	Price    int    `validate:"required" json:"price"`
	Discount *int   `json:"discount"`
}

type ProductResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Barcode  string `json:"barcode"`
	Stock    int    `json:"stock"`
	Price    int    `json:"price"`
	Discount *int   `json:"discount"`
}
