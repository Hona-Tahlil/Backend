package service

import (
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type UserService struct {
	userRepository *postgres.UserRepository
}

func NewGeneralService(userRepository *postgres.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) Login(loginInfo user.LoginRequest) user.LoginResponse {
	token := loginInfo.Email + "/" + loginInfo.Password

	return user.LoginResponse{
		JWTToken: token,
	}
}

func (us *UserService) validateDuplicatePhone(email string) error {
	return nil
}

func (us *UserService) passwordValidation(password string) error {
	return nil
}
func (us *UserService) GenerateFromPassword(password string, cost int) error {

	return nil
}

func (us *UserService) Register(registerInfo user.RegisterRequest) error {
	//validate Duplicate Email
	err := us.validateDuplicatePhone(registerInfo.Email)
	if err != nil {
		return err
	}
	//password Validation
	err = us.passwordValidation(registerInfo.Password)
	if err != nil {
		return err
	}
	//Hash Password

	//Insertion

	return nil
}

func (us *UserService) VerifyEmail(verifyEmailInfo user.VerifyEmailRequest) error {
	return nil
}

func (us *UserService) ForgotPassword(forgetPasswordInfo user.ForgotPasswordRequest) error {
	return nil
}
