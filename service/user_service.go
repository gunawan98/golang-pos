package service

import (
	"context"

	"github.com/gunawan98/golang-restfull-api/model/web"
)

type UserService interface {
	Authenticate(ctx context.Context, request web.LoginRequest) web.UserResponse
	Create(ctx context.Context, request web.UserCreateRequest) web.UserResponse
	Update(ctx context.Context, request web.UserUpdateRequest) web.UserResponse
	Delete(ctx context.Context, userId int)
	FindById(ctx context.Context, userId int) web.UserResponse
	FindAll(ctx context.Context) []web.UserResponse
}
