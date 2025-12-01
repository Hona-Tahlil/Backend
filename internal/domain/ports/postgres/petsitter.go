package domainpostgres

import "hona/backend/internal/domain/entities"

type PetSitterRepository interface {
	CreatePetSitter(petSitter *entities.PetSitter) error
	UpdatePetSitter(petSitter *entities.PetSitter) error
	PreloadServices(petSitter *entities.PetSitter) error
	FindPetSitterByUserID(id uint) (*entities.PetSitter, error)
	GetAllPetSitters(limit, offset int) ([]entities.PetSitter, error)
	GetPetSittersCount() (int64, error)
}
