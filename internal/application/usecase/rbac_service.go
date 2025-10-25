package usecase

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/domain/entities"
)

type RBACService interface {
	GetRoleResponse(role entities.Role) *rbac.RoleResponse
	GetUserInfosResponse(users []entities.User) []rbac.UserInfoResponse
	GetAllRolesWithUsers() ([]rbac.RoleWithUsersResponse, error)
	GetRoleWithUsersByID(info rbac.GetRoleByIDRequest) (*rbac.RoleWithUsersResponse, error)
	GetRoleWithUsersByType(info rbac.GetRoleByTypeRequest) (*rbac.RoleWithUsersResponse, error)
	GetRoleByID(info rbac.GetRoleByIDRequest) (*rbac.RoleResponse, error)
	GetRoleByType(info rbac.GetRoleByTypeRequest) (*rbac.RoleResponse, error)
	GetAllRoles() ([]rbac.RoleResponse, error)
	GetUserRolesByID(info rbac.GetUserRolesByIDRequest) ([]rbac.RoleResponse, error)
	GetUserRolesByEmail(info rbac.GetUserRolesByEmailRequest) ([]rbac.RoleResponse, error)
	RemoveRoleFromUserByID(info rbac.RemoveRoleFromUserByIDRequest) error
	RemoveRoleFromUserByEmail(info rbac.RemoveRoleFromUserByEmailRequest) error
	AddRoleToUserByID(info rbac.AddRoleToUserByIDRequest) error
	AddRoleToUserByEmail(info rbac.AddRoleToUserByEmailRequest) error
	AddRole(info rbac.AddRoleRequest) error
	RemoveRoleByID(info rbac.RemoveRoleByIDRequest) error
	RemoveRoleByType(info rbac.RemoveRoleByTypeRequest) error
	AddPermissionToRole(info rbac.AddPermissionToRoleRequest) error
	RemovePermissionFromRole(info rbac.RemovePermissionFromRoleRequest) error
	GetPermissionRoles(info rbac.GetPermissionRolesRequest) ([]rbac.RoleResponse, error)
}
