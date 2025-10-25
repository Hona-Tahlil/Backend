package domainpostgres

import "hona/backend/internal/domain/entities"

type RBACRepository interface {
	GetRoleByID(roleID uint) (*entities.Role, error)
	GetRoleUsersByID(roleID uint) ([]entities.User, error)
	GetAllRoles() ([]entities.Role, error)
	GetRoleByType(roleType string) (*entities.Role, error)
	RemoveRoleFromUser(user entities.User, role entities.Role) error
	AddRoleToUser(user entities.User, role entities.Role) error
	AddRole(roleType string, description *string) error
	RemoveRole(role entities.Role) error
	AddPermissionToRole(role entities.Role, permission entities.Permission) error
	GetPermissionByID(id uint) (*entities.Permission, error)
	RemovePermissionFromRole(role entities.Role, permission entities.Permission) error
	GetPermissionRolesByID(permissionID uint) ([]entities.Role, error)
}
