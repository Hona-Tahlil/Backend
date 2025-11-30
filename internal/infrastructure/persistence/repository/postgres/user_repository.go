package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (up *UserRepository) FindUserByEmail(email string) (*entities.User, error) {
	var foundUser entities.User

	if result := up.db.First(&foundUser, "email = ?", email); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundUser, nil
}

func (up *UserRepository) FindUserByID(userID uint) (*entities.User, error) {
	var foundUser entities.User

	if result := up.db.First(&foundUser, userID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundUser, nil
}

func (up *UserRepository) CreateUser(user *entities.User) error {
	return up.db.Create(user).Error
}

func (up *UserRepository) DeleteUserByEmail(email string) error {
	return up.db.Where("email = ?", email).Delete(&entities.User{}).Error
}


func (up *UserRepository) SaveUser(user *entities.User) error {
	return up.db.Save(user).Error
}

func (up *UserRepository) GetRoleUsersByID(roleID uint, limit, offset int) ([]entities.User, error) {
	var users []entities.User

	err := up.db.
		Joins("JOIN user_roles ur ON ur.user_id = users.id").
		Where("ur.role_id = ?", roleID).
		Preload("Roles").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (up *UserRepository) PreloadFields(user *entities.User, fields []string) error {
	for _, field := range fields {
		err := up.db.Preload(field).First(user, user.ID).Error
		if err != nil {
			return err
		}
	}
	return nil
}
