package usecase

<<<<<<< HEAD
type PetSitterService interface {
	GetAllPetSitters(page, count int) (interface{}, error)
}
=======
import (
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
)

type PetSitterService interface {
	GetPetSitterFreeSlotsResponse(petSitter *entities.PetSitter) ([]calendarslot.CalendarSlotInfoResponse, error)
	GetServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error)
	GetAvailableServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error)
	GetPetSitterByID(id uint) (*entities.PetSitter, error)
	GetPetSitterByUserID(id uint) (*entities.PetSitter, error)
	PreloadFields(petSitter *entities.PetSitter, fields []string) error
	ValidatePets(pets []entities.Pet, petKinds []enums.PetKind) error
	ValidateService(services []entities.Service, serviceID uint) (*entities.Service, error)
	AutoUpdateSlots(petSitter *entities.PetSitter, calendarSlots []entities.CalendarSlot, accept bool) error
}
>>>>>>> dev
