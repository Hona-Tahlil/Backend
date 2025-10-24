package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	domainjwt "hona/backend/internal/domain/jwt"
	"hona/backend/internal/domain/ports"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	jwtService  domainjwt.JWTService
	unitOfWork  ports.UnitOfWork
	rbacService usecase.RBACService
}

func NewUserService(unitOfWork ports.UnitOfWork, jwtService domainjwt.JWTService, rbacService usecase.RBACService) *UserService {
	return &UserService{
		unitOfWork:  unitOfWork,
		jwtService:  jwtService,
		rbacService: rbacService,
	}
}

func (us *UserService) Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, int, error) {
	foundUser, err := us.findVerifiedUserByEmail(loginInfo.Email)
	if err != nil {
		if _, ok := err.(*exceptions.NotFoundError); ok {
			err = exceptions.NewInvalidCredentialsError("no user found with that email")
		}
		return nil, "", 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginInfo.Password)); err != nil {
		invalidCredentialsErr := exceptions.NewInvalidCredentialsError("password is wrong")
		return nil, "", 0, invalidCredentialsErr
	}

	accessToken, refreshToken, expireTime := us.jwtService.GenerateTokens(foundUser.ID, loginInfo.RememberMe)

	roles := us.rbacService.GetRolesResponse(*foundUser)

	return &user.LoginResponse{
		AccessToken: accessToken,
		Roles:       roles,
	}, refreshToken, expireTime, nil
}

func (us *UserService) findVerifiedUserByEmail(email string) (*entities.User, error) {
	foundUser, err := us.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !foundUser.IsEmailVerified {
		notVerifiedErr := exceptions.NewNotVerifiedError()
		return nil, notVerifiedErr
	}

	return foundUser, nil
}

func (us *UserService) FindUserByEmail(email string) (*entities.User, error) {
	foundUser, err := us.unitOfWork.Factory().UserRepository().FindUserByEmail(email)
	if foundUser == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.Fields.User)
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (us *UserService) findVerifiedUserByID(id uint) (*entities.User, error) {
	foundUser, err := us.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	if !foundUser.IsEmailVerified {
		notVerifiedErr := exceptions.NewNotVerifiedError()
		return nil, notVerifiedErr
	}

	return foundUser, nil
}

func (us *UserService) FindUserByID(id uint) (*entities.User, error) {
	foundUser, err := us.unitOfWork.Factory().UserRepository().FindUserByID(id)
	if foundUser == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.Fields.User)
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return foundUser, nil
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

func (us *UserService) RefreshTokens(refreshTokenInfo rbac.RefreshTokenRequest) (*rbac.RefreshTokenResponse, string, int, error) {
	accessToken, refreshToken, userID, expireTime := us.jwtService.RefreshTokens(refreshTokenInfo.RefreshToken)

	foundUser, err := us.FindUserByID(userID)
	if err != nil {
		return nil, "", 0, err
	}

	roles := us.rbacService.GetRolesResponse(*foundUser)

	return &rbac.RefreshTokenResponse{
		AccessToken: accessToken,
		Roles:       roles,
	}, refreshToken, expireTime, nil
}
