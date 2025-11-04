package domainpostgres

import "hona/backend/internal/domain/entities"

type PetRepository interface {
	CreatePet(pet *entities.Pet) error
	FindPet(name string, userID uint) (*entities.Pet, error)
	FindPetByID(id uint) (*entities.Pet, error)
	UpdatePet(pet *entities.Pet) error
	RemovePet(pet *entities.Pet) error
	PreloadUserPets(user *entities.User) error
}
