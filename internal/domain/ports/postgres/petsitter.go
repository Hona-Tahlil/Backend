package domainpostgres

import "hona/backend/internal/domain/entities"

type PetSitterRepository interface {
	CreatePetSitter(petSitter *entities.PetSitter) error
	SavePetSitter(petSitter *entities.PetSitter) error
	PreloadSchedule(petSitter *entities.PetSitter) error
}
