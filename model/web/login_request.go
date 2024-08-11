package web

type LoginRequest struct {
	Username string `validate:"required,max=50,min=1" json:"username"`
	Password string `validate:"required,max=100,min=1" json:"password"`
}

type RefreshTokenRequest struct {
	Refresh string `validate:"required" json:"refresh"`
}
