package domainpostgres

import "hona/backend/internal/domain/entities"

type UserRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
	DeleteUserByEmail(email string) error
	CreateUser(user *entities.User) error
	SaveUser(user *entities.User) error
}
