package postgres

import (
	"hona/backend/internal/domain/entity"
	"hona/backend/internal/infrastructure/database"
)

type UserRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (up *UserRepository) FindUserByEmail(email string) *entity.User {
	var user entity.User
	if err := up.db.GetDB().First(&user, email); err != nil {
		// TODO: Proper error handling
		panic(err)
	}
	return &user
}
