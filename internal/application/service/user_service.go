package service

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	domainjwt "hona/backend/internal/domain/jwt"
	"hona/backend/internal/domain/ports"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	jwtService domainjwt.JWTService
	unitOfWork ports.UnitOfWork
}

func NewUserService(unitOfWork ports.UnitOfWork, jwtService domainjwt.JWTService) *UserService {
	return &UserService{
		unitOfWork: unitOfWork,
		jwtService: jwtService,
	}
}

func (us *UserService) GetRolesResponse(user entities.User) []rbac.RoleResponse {
	r := make([]rbac.RoleResponse, 0)
	for _, role := range user.Roles {
		p := make([]rbac.PermissionResponse, 0)
		for _, per := range role.Permissions {
			des := ""
			if per.Description != nil {
				des = *per.Description
			}
			p = append(p, rbac.PermissionResponse{
				ID:          per.ID,
				Name:        per.Type.String(),
				Description: des,
				Category:    per.Category.String(),
			})
		}
		des := ""
		if role.Description != nil {
			des = *role.Description
		}
		r = append(r, rbac.RoleResponse{
			ID:          role.ID,
			Name:        role.Type,
			Description: des,
			Permissions: p,
		})
	}
	return r
}

func (us *UserService) Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, int, error) {
	foundUser, err := us.FindUserByEmail(loginInfo.Email)
	fmt.Println("✅ 22")

	if err != nil {
		if _, ok := err.(*exceptions.NotFoundError); !ok {
			return nil, "", 0, err
		}
		fmt.Println("✅ 33")

		invalidCredentialsErr := exceptions.NewInvalidCredentialsError("password is wrong")
		return nil, "", 0, invalidCredentialsErr
	}
	fmt.Println("✅ 44")
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginInfo.Password)); err != nil {
		invalidCredentialsErr := exceptions.NewInvalidCredentialsError("password is wrong")
		return nil, "", 0, invalidCredentialsErr
	}
	fmt.Println("✅ 55")
	accessToken, refreshToken, expireTime := us.jwtService.GenerateTokens(foundUser.ID, loginInfo.RememberMe)
	fmt.Println("✅ 66")
	roles := us.GetRolesResponse(*foundUser)
	r := &user.LoginResponse{
		AccessToken: accessToken,
		Roles:      roles,
	}
	fmt.Println("✅ 77")
	return r, refreshToken, expireTime, nil
}

func (us *UserService) GetUserInfosResponse(users []entities.User) []rbac.UserInfoResponse {
	r := make([]rbac.UserInfoResponse, 0)
	for _, user := range users {
		r = append(r, rbac.UserInfoResponse{
			Email: user.Email,
		})
	}
	return r
}

func (us *UserService) GetRoleUsersByID(roleID uint, limit, offset int) ([]entities.User, error) {
	userRepo := us.unitOfWork.Factory().UserRepository()
	users, err := userRepo.GetRoleUsersByID(roleID, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
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
	userRepo := us.unitOfWork.Factory().UserRepository()
	foundUser, err := userRepo.FindUserByEmail(email)
	if foundUser == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (us *UserService) FindVerifiedUserByID(id uint) (*entities.User, error) {
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
	userRepo := us.unitOfWork.Factory().UserRepository()
	foundUser, err := userRepo.FindUserByID(id)
	if foundUser == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
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

	roles := us.GetRolesResponse(*foundUser)

	return &rbac.RefreshTokenResponse{
		AccessToken: accessToken,
		Roles:       roles,
	}, refreshToken, expireTime, nil
}
