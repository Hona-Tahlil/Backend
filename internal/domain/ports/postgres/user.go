package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type UserRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
	DeleteUserByEmail(email string) error
	CreateUser(user *entities.User) error
	SaveUser(user *entities.User) error
	FindUserByID(userID uint) (*entities.User, error)
	GetRoleUsersByID(roleID uint, options *postgres.QueryOptions) ([]entities.User, error)
	PreloadPetSitter(user *entities.User) error
	PreloadFields(user *entities.User, fields []string) error
	PreloadAddress(user *entities.User) error
}
