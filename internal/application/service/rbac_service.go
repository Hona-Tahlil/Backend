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

func NewRBACService(unitOfWork ports.UnitOfWork, userService usecase.UserService) *RBACService {
	return &RBACService{
		unitOfWork:  unitOfWork,
		userService: userService,
	}
}

func (rs *RBACService) GetRoleResponse(role entities.Role) *rbac.RoleResponse {
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
	r := &rbac.RoleResponse{
		ID:          role.ID,
		Name:        role.Type,
		Description: des,
		Permissions: p,
	}
	return r
}

func (rs *RBACService) ListRolesWithUsers(info rbac.ListRolesWithUsersRequest) ([]rbac.RoleWithUsersResponse, error) {
	r := make([]rbac.RoleWithUsersResponse, 0)

	limit := info.Count
	offset := (info.Page - 1) * limit

	roles, err := rs.unitOfWork.Factory().RBACRepository().GetAllRoles()
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		res, err := rs.getRoleWithUsers(&role, limit, offset)
		if err != nil {
			return nil, err
		}
		r = append(r, *res)
	}

	return r, nil
}

func (rs *RBACService) getRoleWithUsers(role *entities.Role, limit, offset int) (*rbac.RoleWithUsersResponse, error) {
	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	users, err := rs.userService.GetRoleUsersByID(role.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	if err := rbacRepo.PreloadRolePermissions(role); err != nil {
		return nil, err
	}

	return &rbac.RoleWithUsersResponse{
		Role:  *rs.GetRoleResponse(*role),
		Users: rs.userService.GetUserInfosResponse(users),
	}, nil
}

func (rs *RBACService) findRoleByID(roleID uint) (*entities.Role, error) {
	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	foundRole, err := rbacRepo.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}

	if foundRole == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Role)
		return nil, NotFoundError
	}

	return foundRole, nil
}

func (rs *RBACService) findRoleByType(roleType string) (*entities.Role, error) {
	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	foundRole, err := rbacRepo.GetRoleByType(roleType)
	if err != nil {
		return nil, err
	}

	if foundRole == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Role)
		return nil, NotFoundError
	}

	return foundRole, nil
}

func (rs *RBACService) findPermissionByID(permissionID uint) (*entities.Permission, error) {
	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	foundPermission, err := rbacRepo.GetPermissionByID(permissionID)
	if err != nil {
		return nil, err
	}

	if foundPermission == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Permission)
		return nil, NotFoundError
	}

	return foundPermission, nil
}

func (rs *RBACService) GetRoleWithUsersByID(info rbac.GetRoleWithUsersByIDRequest) (*rbac.RoleWithUsersResponse, error) {
	role, err := rs.findRoleByID(info.ID)
	if err != nil {
		return nil, err
	}

	limit := info.Count
	offset := (info.Page - 1) * limit

	return rs.getRoleWithUsers(role, limit, offset)
}

func (rs *RBACService) GetRoleWithUsersByType(info rbac.GetRoleWithUsersByTypeRequest) (*rbac.RoleWithUsersResponse, error) {
	role, err := rs.findRoleByType(info.Type)
	if err != nil {
		return nil, err
	}

	limit := info.Count
	offset := (info.Page - 1) * limit

	return rs.getRoleWithUsers(role, limit, offset)
}

func (rs *RBACService) GetRoleByID(info rbac.GetRoleByIDRequest) (*rbac.RoleResponse, error) {
	role, err := rs.findRoleByID(info.ID)
	if err != nil {
		return nil, err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	if err := rbacRepo.PreloadRolePermissions(role); err != nil {
		return nil, err
	}

	return rs.GetRoleResponse(*role), nil
}

func (rs *RBACService) GetRoleByType(info rbac.GetRoleByTypeRequest) (*rbac.RoleResponse, error) {
	role, err := rs.findRoleByType(info.Type)
	if err != nil {
		return nil, err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	if err := rbacRepo.PreloadRolePermissions(role); err != nil {
		return nil, err
	}

	return rs.GetRoleResponse(*role), nil
}

func (rs *RBACService) GetAllRoles() ([]rbac.RoleResponse, error) {
	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	r := make([]rbac.RoleResponse, 0)

	roles, err := rbacRepo.GetAllRoles()
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if err := rbacRepo.PreloadRolePermissions(&role); err != nil {
			return nil, err
		}
		res := rs.GetRoleResponse(role)
		r = append(r, *res)
	}

	return r, nil
}

func (rs *RBACService) GetUserRolesByID(info rbac.GetUserRolesByIDRequest) ([]rbac.RoleResponse, error) {
	user, err := rs.userService.FindUserByID(info.ID, true)
	if err != nil {
		return nil, err
	}

	return rs.userService.GetRolesResponse(*user), nil
}

func (rs *RBACService) GetUserRolesByEmail(info rbac.GetUserRolesByEmailRequest) ([]rbac.RoleResponse, error) {
	user, err := rs.userService.FindUserByEmail(info.Email, true)
	if err != nil {
		return nil, err
	}

	return rs.userService.GetRolesResponse(*user), nil
}

func (rs *RBACService) RemoveRoleFromUserByID(info rbac.RemoveRoleFromUserByIDRequest) error {
	user, err := rs.userService.FindUserByID(info.UserID, false)
	if err != nil {
		return err
	}

	role, err := rs.findRoleByID(info.RoleID)
	if err != nil {
		return err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	err = rbacRepo.RemoveRoleFromUser(user, role)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RBACService) RemoveRoleFromUserByEmail(info rbac.RemoveRoleFromUserByEmailRequest) error {
	user, err := rs.userService.FindUserByEmail(info.UserEmail, false)
	if err != nil {
		return err
	}

	role, err := rs.findRoleByID(info.RoleID)
	if err != nil {
		return err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	err = rbacRepo.RemoveRoleFromUser(user, role)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RBACService) AddRoleToUserByID(info rbac.AddRoleToUserByIDRequest) error {
	user, err := rs.userService.FindUserByID(info.UserID, false)
	if err != nil {
		return err
	}

	role, err := rs.findRoleByID(info.RoleID)
	if err != nil {
		return err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	err = rbacRepo.AddRoleToUser(user, role)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RBACService) AddRoleToUserByEmail(info rbac.AddRoleToUserByEmailRequest) error {
	user, err := rs.userService.FindUserByEmail(info.UserEmail, false)
	if err != nil {
		return err
	}

	role, err := rs.findRoleByID(info.RoleID)
	if err != nil {
		return err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	err = rbacRepo.AddRoleToUser(user, role)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RBACService) AddRole(info rbac.AddRoleRequest) error {
	if info.Description != nil && strings.TrimSpace(*info.Description) == "" {
		info.Description = nil
	}
	_, err := rs.findRoleByType(info.Type)
	if err == nil {
		var ce exceptions.ConflictErrors
		ce.Add(bootstrap.Run().Constants.ErrorFields.Role, bootstrap.Run().Constants.ErrorTags.AlreadyExist)
		return &ce
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	if err := rbacRepo.AddRole(info.Type, info.Description); err != nil {
		return err
	}
	return nil
}

func (rs *RBACService) RemoveRoleByID(info rbac.RemoveRoleByIDRequest) error {
	role, err := rs.findRoleByID(info.ID)
	if err != nil {
		return err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	err = rbacRepo.RemoveRole(role)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RBACService) RemoveRoleByType(info rbac.RemoveRoleByTypeRequest) error {
	role, err := rs.findRoleByType(info.Type)
	if err != nil {
		return err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	err = rbacRepo.RemoveRole(role)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RBACService) AddPermissionToRole(info rbac.AddPermissionToRoleRequest) error {
	role, err := rs.findRoleByID(info.RoleID)
	if err != nil {
		return err
	}

	permission, err := rs.findPermissionByID(info.PermissionID)
	if err != nil {
		return err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	err = rbacRepo.AddPermissionToRole(role, permission)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RBACService) RemovePermissionFromRole(info rbac.RemovePermissionFromRoleRequest) error {
	role, err := rs.findRoleByID(info.RoleID)
	if err != nil {
		return err
	}

	permission, err := rs.findPermissionByID(info.PermissionID)
	if err != nil {
		return err
	}

	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	err = rbacRepo.RemovePermissionFromRole(role, permission)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RBACService) GetPermissionRoles(info rbac.GetPermissionRolesRequest) ([]rbac.RoleResponse, error) {
	rbacRepo := rs.unitOfWork.Factory().RBACRepository()
	r := make([]rbac.RoleResponse, 0)

	roles, err := rbacRepo.GetPermissionRolesByID(info.PermissionID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if err := rbacRepo.PreloadRolePermissions(&role); err != nil {
			return nil, err
		}
		r = append(r, *rs.GetRoleResponse(role))
	}

	return r, nil
}
