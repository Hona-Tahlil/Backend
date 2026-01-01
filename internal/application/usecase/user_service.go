package usecase

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type UserService interface {
	Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, int, error)
	FindVerifiedUserByEmail(email string) (*entities.User, error)
	FindUserByEmail(email string) (*entities.User, error)
	FindVerifiedUserByID(id uint) (*entities.User, error)
	FindUserByID(id uint) (*entities.User, error)
	RefreshTokens(refreshTokenInfo rbac.RefreshTokenRequest) (*rbac.RefreshTokenResponse, string, int, error)
	GetRolesResponse(user *entities.User) []rbac.RoleResponse
	GetRoleUsersByID(roleID uint, options *postgres.QueryOptions) ([]entities.User, int64, error)
	GetUserInfosResponse(users []entities.User) ([]rbac.UserResponse, error)
	GetUserInfoResponse(userEntity *entities.User) (*rbac.UserResponse, error)
	GetProfile(info user.GetProfileRequest) (*user.ProfileResponse, error)
	GetIdentity(info user.GetIdentityRequest) (*user.IdentityResponse, error)
	UpdateProfile(info user.UpdateProfileRequest) error
	Register(registerInfo user.RegisterRequest) error
	ResetPassword(resetPasswordInfo user.ResetPasswordRequest) error
	ForgotPassword(forgetPasswordInfo user.ForgotPasswordRequest) error
	VerifyEmail(info user.VerifyEmailRequest) error
	SendVerificationEmail(info user.SendVerificationEmailRequest) error
	PreloadFields(user *entities.User, fields []string) error
}
