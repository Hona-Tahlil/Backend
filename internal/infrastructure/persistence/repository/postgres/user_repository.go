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
	var user entities.User
	if err := up.db.First(&user, "email = ?", email); err.Error != nil {
		return nil, err.Error
	}
	return &user, nil
}
  

func (up *UserRepository) DeleteUserByEmail(email string) error {
	return  nil
}

func (up *UserRepository) CreateUser(*entities.User) error {
	return nil
}