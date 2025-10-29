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

	if err := up.db.First(&foundUser, "email = ?", email); err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err.Error
	}
	return &foundUser, nil
}

func (up *UserRepository) CreateUser(user *entities.User) error {
	return up.db.Create(user).Error
}

func (up *UserRepository) DeleteUserByEmail(email string) error {
	return up.db.Where("email = ?", email).Delete(&entities.User{}).Error
}
