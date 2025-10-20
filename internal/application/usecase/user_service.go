package usecase

import "hona/backend/internal/application/dto/user"

type UserService interface {
	Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, int, error)
	Register(registerInfo user.RegisterRequest) error
}
