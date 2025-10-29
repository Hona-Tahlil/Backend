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
