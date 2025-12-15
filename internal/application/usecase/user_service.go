package usecase

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
)

type UserService interface {
	Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, int, error)
	// findVerifiedUserByEmail(email string) (*entities.User, error)
	// findVerifiedUserByID(id uint) (*entities.User, error)
	FindUserByID(id uint) (*entities.User, error)
	FindUserByEmail(email string) (*entities.User, error)
	// validateDuplicatePhone(email string) error
	// passwordValidation(password string) error
	// GenerateFromPassword(password string, cost int) error
	// Register(registerInfo user.RegisterRequest) error
	// VerifyEmail(verifyEmailInfo user.VerifyEmailRequest) error
	// ForgotPassword(forgetPasswordInfo user.ForgotPasswordRequest) error
	RefreshTokens(refreshTokenInfo rbac.RefreshTokenRequest) (*rbac.RefreshTokenResponse, string, int, error)
	GetRolesResponse(user entities.User) []rbac.RoleResponse
	GetRoleUsersByID(roleID uint, limit, offset int) ([]entities.User, error)
	GetUserInfosResponse(users []entities.User) []rbac.UserInfoResponse
	Register(registerInfo user.RegisterRequest) error
	ResetPassword(resetPasswordInfo user.ResetPasswordRequest) error
	ForgotPassword(forgetPasswordInfo user.ForgotPasswordRequest) error
	VerifyEmail(info user.VerifyEmailRequest) error
	SendVerificationEmail(info user.SendVerificationEmailRequest) error
}
