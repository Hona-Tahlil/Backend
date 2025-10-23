package service

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/infrastructure/jwt"
	"hona/backend/internal/infrastructure/persistence"

	"golang.org/x/crypto/bcrypt"
)

// TODO: interface for service and repo

type UserService struct {
	jwtService  *jwt.JWTService
	unitOfWork  *persistence.UnitOfWork
	rbacService *RBACService
}

func NewUserService(unitOfWork *persistence.UnitOfWork, jwtService *jwt.JWTService, rbacService *RBACService) *UserService {
	return &UserService{
		unitOfWork:  unitOfWork,
		jwtService:  jwtService,
		rbacService: rbacService,
	}
}

func (us *UserService) Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, int, error) {
	foundUser, err := us.findVerifiedUserByEmail(loginInfo.Email)
	if err != nil {
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
	foundUser, err := us.unitOfWork.Factory().UserRepository().FindUserByEmail(email)
	if foundUser == nil {
		invalidCredentialsErr := exceptions.NewInvalidCredentialsError("email not found")
		return nil, invalidCredentialsErr
	}

	if err != nil {
		return nil, err
	}

	if !foundUser.IsVerified {
		notVerifiedErr := exceptions.NewNotVerifiedError()
		return nil, notVerifiedErr
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

func (us *UserService) RefreshTokens(refreshTokenInfo rbac.RefreshTokenRequest) (*rbac.RefreshTokenResponse, string, error) {
	accessToken, refreshToken, userID := us.jwtService.RefreshTokens(refreshTokenInfo.RefreshToken)

	foundUser, err := us.unitOfWork.Factory().UserRepository().FindUserByID(userID)
	if err != nil {
		unauthorizedError := exceptions.NewUnauthorizedError("user not found")
		return nil, "", unauthorizedError
	}

	p := make([]rbac.PermissionResponse, 0)
	for _, role := range foundUser.Roles {
		for _, per := range role.Permissions {
			p = append(p, rbac.PermissionResponse{
				ID:   per.ID,
				Name: per.Type.String(),
			})
		}
	}

	return &rbac.RefreshTokenResponse{
		AccessToken: accessToken,
		Permissions: p,
	}, refreshToken, nil
}
