package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type PetSitterRepository interface {
	CreatePetSitter(petSitter *entities.PetSitter) error
	PreloadServices(petSitter *entities.PetSitter) error
	FindPetSitterByUserID(id uint) (*entities.PetSitter, error)
	GetAllPetSitters(limit, offset int) ([]entities.PetSitter, error)
	GetPetSittersCount() (int64, error)
	PreloadFields(petSitter *entities.PetSitter, fields []string) error
	UpdatePetSitter(petSitter *entities.PetSitter) error
	ReplaceSchedule(petSitter *entities.PetSitter, schedule []entities.CalendarSlot) error
	FindPetSitterByID(id uint) (*entities.PetSitter, error)
	FindServiceByID(id uint) (*entities.Service, error)
	CreateService(service *entities.Service) error
	UpdateService(service *entities.Service) error
	DeleteService(service *entities.Service) error
	SearchPetSitters(options *postgres.QueryOptions) ([]*entities.PetSitter, int64, error)
}
