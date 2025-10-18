package usecase

import "hona/backend/internal/application/dto/user"

type UserService interface {
	Login(loginInfo user.LoginRequest) user.LoginResponse
	Register(registerInfo user.RegisterRequest) error
}
