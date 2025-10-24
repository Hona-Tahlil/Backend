package domainpostgres

import "hona/backend/internal/domain/entities"

type RBACRepository interface {
	GetRoleByID(roleID uint) (*entities.Role, error)
	GetRoleUsersByID(roleID uint) ([]entities.User, error)
	GetAllRoles() ([]entities.Role, error)
	GetRoleByType(roleType string) (*entities.Role, error)
	RemoveRoleFromUserByID(user entities.User, roleID uint) error
	AddRoleToUserByID(user entities.User, roleID uint) error
	AddRole(roleType string, description *string) error
}
