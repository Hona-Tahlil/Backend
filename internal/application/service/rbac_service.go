package service

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/domain/entities"
)

type RBACService struct {
}

func NewRBACService() *RBACService {
	return &RBACService{}
}

func (rs *RBACService) GetRolesResponse(user entities.User) []rbac.RoleResponse {
	r := make([]rbac.RoleResponse, 0)
	for _, role := range user.Roles {
		p := make([]rbac.PermissionResponse, 0)
		for _, per := range role.Permissions {
			p = append(p, rbac.PermissionResponse{
				ID:          per.ID,
				Name:        per.Type.String(),
				Description: *per.Description,
				Category:    per.Category.String(),
			})
		}
		r = append(r, rbac.RoleResponse{
			ID:          role.ID,
			Name:        role.Type,
			Description: *role.Description,
			Permissions: p,
		})
	}
	return r
}

// TODO: response
func (rs *RBACService) GetAllRolesWithUsers() {

}
