package domainpostgres

import "hona/backend/internal/domain/entities"

type PetSitterRepository interface {
	PreloadFields(petSitter *entities.PetSitter, fields []string) error
}
