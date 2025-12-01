package domainpostgres

import "hona/backend/internal/domain/entities"

type UserRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
	DeleteUserByEmail(email string) error
	CreateUser(user *entities.User) error
	SaveUser(user *entities.User) error
	FindUserByID(userID uint) (*entities.User, error)
	GetRoleUsersByID(roleID uint, limit, offset int) ([]entities.User, error)
<<<<<<< HEAD
	PreloadPetSitter(user *entities.User) error
	UpdateUser(user *entities.User) error
	DeleteUserByEmail(email string) error
	CreateUser(user *entities.User) error	
	PreloadAddress(user *entities.User) error
	
=======
	PreloadFields(user *entities.User, fields []string) error
>>>>>>> dev
}
