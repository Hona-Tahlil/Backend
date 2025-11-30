package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/dsl"
)

type UserRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
	DeleteUserByEmail(email string) error
	CreateUser(user *entities.User) error
	SaveUser(user *entities.User) error
	FindUserByID(userID uint) (*entities.User, error)
	GetRoleUsersByID(roleID uint,queryoptins *dsl.ParsedQuery) ([]entities.User, error)
}
