package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/dsl"
)

type PetSitterRepository interface {
	CreatePetSitter(petSitter *entities.PetSitter) error
	UpdatePetSitter(petSitter *entities.PetSitter) error
	FindPetSitterByUserID(userID uint) (*entities.PetSitter, error)
	PreloadServices(petSitter *entities.PetSitter) error
	SearchPetSitters(offset int, limit int, filters []dsl.Filter, sorts []dsl.Sort) ([]*entities.PetSitter, int64, error)
}
