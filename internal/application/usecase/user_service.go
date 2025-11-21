package usecase

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
)

type UserService interface {
	Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, int, error)
	FindVerifiedUserByEmail(email string) (*entities.User, error)
	FindUserByEmail(email string) (*entities.User, error)
	FindVerifiedUserByID(id uint) (*entities.User, error)
	FindUserByID(id uint) (*entities.User, error)
	RefreshTokens(refreshTokenInfo rbac.RefreshTokenRequest) (*rbac.RefreshTokenResponse, string, int, error)
	GetRolesResponse(user entities.User) []rbac.RoleResponse
	GetRoleUsersByID(roleID uint, limit, offset int) ([]entities.User, error)
	GetUserInfosResponse(users []entities.User) []rbac.UserInfoResponse
	PreloadFields(user *entities.User, fields []string) error
}
