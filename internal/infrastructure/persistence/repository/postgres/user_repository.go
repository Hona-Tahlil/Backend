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

func (up *UserRepository) GetRoleUsersByID(roleID uint, options *QueryOptions) ([]entities.User, int64, error) {
	var users []entities.User

	baseQuery := up.db.Model(&entities.User{}).
		Joins("JOIN user_roles ur ON ur.user_id = users.id").
		Where("ur.role_id = ?", roleID)

	newDB, total := ApplyModifiers(baseQuery, *options)

	err := newDB.
		Preload("Roles.Permissions").
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
func (up *UserRepository) PreloadPetSitter(user *entities.User) error {
	return up.db.Preload("PetSitter").First(user, user.ID).Error
}

func (up *UserRepository) PreloadAddress(user *entities.User) error {
	return up.db.Preload("Address").First(user, user.ID).Error
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
