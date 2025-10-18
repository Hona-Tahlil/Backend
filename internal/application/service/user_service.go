package service

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/infrastructure/jwt"
	"hona/backend/internal/infrastructure/persistence"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	jwtService *jwt.JWTService
	unitOfWork *persistence.UnitOfWork
}

func NewUserService(unitOfWork *persistence.UnitOfWork, jwtService *jwt.JWTService) *UserService {
	return &UserService{
		unitOfWork: unitOfWork,
		jwtService: jwtService,
	}
}

func (us *UserService) Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, error) {

	foundUser, err := us.unitOfWork.Factory().UserRepository().FindUserByEmail(loginInfo.Email)
	if err != nil {
		invalidCredentialsErr := &exceptions.AuthError{
			Type: "INVALID_CREDENTIALS",
		}
		return nil, "", invalidCredentialsErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginInfo.Password)); err != nil {
		invalidCredentialsErr := &exceptions.AuthError{
			Type: "INVALID_CREDENTIALS",
		}
		return nil, "", invalidCredentialsErr
	}

	accessToken, refreshToken := us.jwtService.GenerateTokens(foundUser.ID, loginInfo.RememberMe)

	p := make([]rbac.PermissionResponse, 0)
	for _, role := range foundUser.Roles {
		for _, per := range role.Permissions {
			p = append(p, rbac.PermissionResponse{
				ID:   per.ID,
				Name: per.Type.String(),
			})
		}
	}

	return &user.LoginResponse{
		AccessToken: accessToken,
		Permissions: p,
	}, refreshToken, nil
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
