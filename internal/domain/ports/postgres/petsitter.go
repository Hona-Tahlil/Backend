package domainpostgres

import "hona/backend/internal/domain/entities"

type PetSitterRepository interface {
	CreatePetSitter(petSitter *entities.PetSitter) error
	UpdatePetSitter(petSitter *entities.PetSitter) error
	PreloadServices(petSitter *entities.PetSitter) error
	FindPetSitterByID(id uint) (*entities.PetSitter, error)
}
