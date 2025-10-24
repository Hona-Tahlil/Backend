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

func (rr *RBACRepository) GetRoleUsersByID(roleID uint) ([]entities.User, error) {
	var users []entities.User
	err := rr.db.
		Joins("JOIN user_roles ur ON ur.user_id = users.id").
		Where("ur.role_id = ?", roleID).
		Preload("Roles").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (rr *RBACRepository) GetAllRoles() ([]entities.Role, error) {
	var roles []entities.Role
	if err := rr.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (rr *RBACRepository) RemoveRoleFromUserByID(user entities.User, roleID uint) error {
	role, err := rr.GetRoleByID(roleID)
	if err != nil {
		return err
	}
	if err := rr.db.Model(&user).Association("Roles").Delete(role); err != nil {
		panic(err)
	}
	return nil
}

func (rr *RBACRepository) AddRoleToUserByID(user entities.User, roleID uint) error {
	role, err := rr.GetRoleByID(roleID)
	if err != nil {
		return err
	}
	if err := rr.db.Model(&user).Association("Roles").Append(&role); err != nil {
		panic(err)
	}
	return nil
}

func (rr *RBACRepository) AddRole(roleType string, description *string) error {
	role := entities.Role{
		Type:        roleType,
		Description: description,
	}
	if err := rr.db.Create(&role).Error; err != nil {
		return err
	}
	return nil
}
