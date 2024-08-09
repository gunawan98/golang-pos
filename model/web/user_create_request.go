package web

type UserCreateRequest struct {
	Username string `validate:"required,max=50,min=1" json:"username"`
	Password string `validate:"required,max=100,min=1" json:"password"`
	Role     string `validate:"required,role" json:"role"`
}
