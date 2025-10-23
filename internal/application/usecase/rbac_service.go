package usecase

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/domain/entities"
)

type RBACService interface {
	GetRolesResponse(user entities.User) []rbac.RoleResponse
}
