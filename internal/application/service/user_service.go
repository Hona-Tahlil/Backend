package service

import (
	"hona/backend/internal/application/dto/login"
	"hona/backend/internal/infrastructure/repository/postgres"
)

type GeneralService struct {
	userRepository *postgres.UserRepository
}

func NewGeneralService(userRepository *postgres.UserRepository) *GeneralService {
	return &GeneralService{
		userRepository: userRepository,
	}
}

func (gs *GeneralService) Login(loginInfo login.LoginRequest) login.LoginResponse {
	// TODO: actually implement
	token := loginInfo.Email + "/" + loginInfo.Password

	// TODO: call repo?

	return login.LoginResponse{
		JWTToken: token,
	}
}
