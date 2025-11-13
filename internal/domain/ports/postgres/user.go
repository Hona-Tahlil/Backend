package domainpostgres

import "hona/backend/internal/domain/entities"

type UserRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
	FindUserByID(userID uint) (*entities.User, error)
	GetRoleUsersByID(roleID uint, limit, offset int) ([]entities.User, error)
	PreloadPets(user *entities.User) error
	PreloadPetSitter(user *entities.User) error
}
