package usecase

import (
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/domain/entities"
)

type PetsitterService interface {
	SearchPetSitters(info petsitter.SearchPetSittersRequest) ([]*entities.PetSitter, int64, error)
}
