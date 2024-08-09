package web

type UserUpdateRequest struct {
	Id       int    `validate:"required"`
	Password string `validate:"required,max=100,min=1" json:"password"`
	Role     string `validate:"required,role" json:"role"`
}
