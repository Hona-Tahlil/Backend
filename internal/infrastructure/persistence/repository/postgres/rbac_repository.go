package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type RBACRepository struct {
	db *gorm.DB
}

func NewRBACRepository(db *gorm.DB) *RBACRepository {
	return &RBACRepository{
		db: db,
	}
}

func (rr *RBACRepository) GetRoleByID(roleID uint) (*entities.Role, error) {
	var foundRole entities.Role

	if result := rr.db.First(&foundRole, roleID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundRole, nil
}

func (rr *RBACRepository) GetRoleByType(roleType string) (*entities.Role, error) {
	var foundRole entities.Role

	if result := rr.db.First(&foundRole, "type = ?", roleType); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundRole, nil
}

func (rr *RBACRepository) GetAllRoles() ([]entities.Role, error) {
	var roles []entities.Role
	if err := rr.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (rr *RBACRepository) RemoveRoleFromUser(user *entities.User, role *entities.Role) error {
	if err := rr.db.Model(user).Association("Roles").Delete(role); err != nil {
		return err
	}

	if err := rr.db.Preload("Roles").First(&user).Error; err != nil {
		return err
	}

	return nil
}

func (rr *RBACRepository) AddRoleToUser(user *entities.User, role *entities.Role) error {
	if err := rr.db.Model(user).Association("Roles").Append(role); err != nil {
		return err
	}

	if err := rr.db.Preload("Roles").First(&user).Error; err != nil {
		return err
	}

	return nil
}

func (rr *RBACRepository) AddRole(roleType string, description *string) error {
	role := entities.Role{
		Type:        roleType,
		Description: description,
		Permissions: make([]entities.Permission, 0),
	}
	if err := rr.db.Create(&role).Error; err != nil {
		return err
	}
	return nil
}

func (rr *RBACRepository) RemoveRole(role *entities.Role) error {
	return rr.db.Delete(&role).Error
}

func (rr *RBACRepository) AddPermissionToRole(role *entities.Role, permission *entities.Permission) error {
	return rr.db.Model(role).Association("Permissions").Append(permission)
}

func (rr *RBACRepository) GetPermissionByID(id uint) (*entities.Permission, error) {
	var foundPermission entities.Permission

	if result := rr.db.First(&foundPermission, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundPermission, nil
}

func (rr *RBACRepository) RemovePermissionFromRole(role *entities.Role, permission *entities.Permission) error {
	return rr.db.Model(&role).Association("Permissions").Delete(&permission)
}

func (rr *RBACRepository) GetPermissionRolesByID(permissionID uint) ([]entities.Role, error) {
	var roles []entities.Role
	err := rr.db.
		Joins("JOIN role_permission rp ON rp.role_id = roles.id").
		Where("rp.permission_id = ?", permissionID).
		Preload("Permissions").
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (rr *RBACRepository) PreloadUserRoles(user *entities.User) error {
	return rr.db.Preload("Roles.Permissions").First(user, user.ID).Error
}

func (rr *RBACRepository) PreloadRolePermissions(role *entities.Role) error {
	return rr.db.Preload("Permissions").First(role, role.ID).Error
}
