package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	"strings"
)

type RBACService struct {
	unitOfWork  ports.UnitOfWork
	userService usecase.UserService
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

func (rs *RBACService) GetRoleResponse(role entities.Role) *rbac.RoleResponse {
	p := make([]rbac.PermissionResponse, 0)
	for _, per := range role.Permissions {
		p = append(p, rbac.PermissionResponse{
			ID:          per.ID,
			Name:        per.Type.String(),
			Description: *per.Description,
			Category:    per.Category.String(),
		})
	}
	r := &rbac.RoleResponse{
		ID:          role.ID,
		Name:        role.Type,
		Description: *role.Description,
		Permissions: p,
	}
	return r
}

func (rs *RBACService) GetUserInfosResponse(users []entities.User) []rbac.UserInfoResponse {
	r := make([]rbac.UserInfoResponse, 0)
	for _, user := range users {
		r = append(r, rbac.UserInfoResponse{
			Username: user.Username,
		})
	}
	return r
}

func (rs *RBACService) GetAllRolesWithUsers() ([]rbac.RoleWithUsersResponse, error) {
	r := make([]rbac.RoleWithUsersResponse, 0)

	roles, err := rs.unitOfWork.Factory().RBACRepository().GetAllRoles()
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		res, err := rs.getRoleWithUsers(&role)
		if err != nil {
			return nil, err
		}
		r = append(r, *res)
	}

	return r, nil
}

func (rs *RBACService) getRoleWithUsers(role *entities.Role) (*rbac.RoleWithUsersResponse, error) {
	users, err := rs.unitOfWork.Factory().RBACRepository().GetRoleUsersByID(role.ID)
	if err != nil {
		return nil, err
	}

	return &rbac.RoleWithUsersResponse{
		Role:  *rs.GetRoleResponse(*role),
		Users: rs.GetUserInfosResponse(users),
	}, nil
}

func (rs *RBACService) findRoleByID(roleID uint) (*entities.Role, error) {
	foundRole, err := rs.unitOfWork.Factory().RBACRepository().GetRoleByID(roleID)
	if foundRole == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.Fields.Role)
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return foundRole, nil
}

func (rs *RBACService) findRoleByType(roleType string) (*entities.Role, error) {
	foundRole, err := rs.unitOfWork.Factory().RBACRepository().GetRoleByType(roleType)
	if foundRole == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.Fields.Role)
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return foundRole, nil
}

func (rs *RBACService) GetRoleWithUsersByID(info rbac.GetRoleByIDRequest) (*rbac.RoleWithUsersResponse, error) {
	role, err := rs.findRoleByID(info.ID)
	if err != nil {
		return nil, err
	}

	return rs.getRoleWithUsers(role)
}

func (rs *RBACService) GetRoleWithUsersByType(info rbac.GetRoleByTypeRequest) (*rbac.RoleWithUsersResponse, error) {
	role, err := rs.findRoleByType(info.Type)
	if err != nil {
		return nil, err
	}

	return rs.getRoleWithUsers(role)
}

func (rs *RBACService) GetRoleByID(info rbac.GetRoleByIDRequest) (*rbac.RoleResponse, error) {
	role, err := rs.findRoleByID(info.ID)
	if err != nil {
		return nil, err
	}

	return rs.GetRoleResponse(*role), nil
}

func (rs *RBACService) GetRoleByType(info rbac.GetRoleByTypeRequest) (*rbac.RoleResponse, error) {
	role, err := rs.findRoleByType(info.Type)
	if err != nil {
		return nil, err
	}

	return rs.GetRoleResponse(*role), nil
}

func (rs *RBACService) GetAllRoles() ([]rbac.RoleResponse, error) {
	r := make([]rbac.RoleResponse, 0)

	roles, err := rs.unitOfWork.Factory().RBACRepository().GetAllRoles()
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		res := rs.GetRoleResponse(role)
		r = append(r, *res)
	}

	return r, nil
}

func (rs *RBACService) GetUserRolesByID(info rbac.GetUserRolesByIDRequest) ([]rbac.RoleResponse, error) {
	user, err := rs.userService.FindUserByID(info.ID)
	if err != nil {
		return nil, err
	}

	return rs.GetRolesResponse(*user), nil
}

func (rs *RBACService) GetUserRolesByEmail(info rbac.GetUserRolesByEmailRequest) ([]rbac.RoleResponse, error) {
	user, err := rs.userService.FindUserByEmail(info.Email)
	if err != nil {
		return nil, err
	}

	return rs.GetRolesResponse(*user), nil
}

func (rs *RBACService) RemoveRoleFromUserByID(info rbac.RemoveRoleFromUserByIDRequest) error {
	user, err := rs.userService.FindUserByID(info.UserID)
	if err != nil {
		return err
	}
	err = rs.unitOfWork.Factory().RBACRepository().RemoveRoleFromUserByID(*user, info.RoleID)
	if err != nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.Fields.Role)
		return NotFoundError
	}
	return nil
}

func (rs *RBACService) RemoveRoleFromUserByEmail(info rbac.RemoveRoleFromUserByEmailRequest) error {
	user, err := rs.userService.FindUserByEmail(info.UserEmail)
	if err != nil {
		return err
	}
	err = rs.unitOfWork.Factory().RBACRepository().RemoveRoleFromUserByID(*user, info.RoleID)
	if err != nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.Fields.Role)
		return NotFoundError
	}
	return nil
}

func (rs *RBACService) AddRoleToUserByID(info rbac.AddRoleToUserByIDRequest) error {
	user, err := rs.userService.FindUserByID(info.UserID)
	if err != nil {
		return err
	}
	err = rs.unitOfWork.Factory().RBACRepository().AddRoleToUserByID(*user, info.RoleID)
	if err != nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.Fields.Role)
		return NotFoundError
	}
	return nil
}

func (rs *RBACService) AddRoleToUserByEmail(info rbac.AddRoleToUserByEmailRequest) error {
	user, err := rs.userService.FindUserByEmail(info.UserEmail)
	if err != nil {
		return err
	}
	err = rs.unitOfWork.Factory().RBACRepository().AddRoleToUserByID(*user, info.RoleID)
	if err != nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.Fields.Role)
		return NotFoundError
	}
	return nil
}

func (rs *RBACService) AddRole(info rbac.AddRoleRequest) error {
	if info.Description != nil && strings.TrimSpace(*info.Description) == "" {
		info.Description = nil
	}
	_, err := rs.findRoleByType(info.Type)
	if err == nil {
		var ce *exceptions.ConflictErrors
		ce.Add(bootstrap.Run().Constants.ErrorFields.Role, bootstrap.Run().Constants.ErrorTags.AlreadyExist)
		return ce
	}
	if err := rs.unitOfWork.Factory().RBACRepository().AddRole(info.Type, info.Description); err != nil {
		return err
	}

	return nil
}
