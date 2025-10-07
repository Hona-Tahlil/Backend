package service

import "hona/backend/internal/application/dto/login"

type GeneralService struct {
}

func NewGeneralService() *GeneralService {
	return &GeneralService{}
}

func (gs *GeneralService) Login(loginInfo login.LoginRequest) login.LoginResponse {
	// TODO: actually implement
	token := loginInfo.Email + "/" + loginInfo.Password

	// TODO: call repo?

	return login.LoginResponse{
		JWTToken: token,
	}
}
