package domainpostgres

import "hona/backend/internal/domain/entities"

type PetSitterRepository interface {
	PreloadSchedule(petSitter *entities.PetSitter) error
}
