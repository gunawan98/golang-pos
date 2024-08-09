package web

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
