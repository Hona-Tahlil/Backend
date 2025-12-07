package usecase

import (
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/petsitter"
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
	CreateSignupSession(PetsitterInfo petsitter.GetPetSitterRequest) (*petsitter.PetSitterStatusResponse, error)
	GetPersonalInfo(userID uint) (*petsitter.PersonalInfoResponse, error)
	UploadDocuments(info petsitter.UploadDocumentsRequest) error
	GetDocuments(userID uint) (*petsitter.DocumentResponse, error)
	SubmitSkills(SkillsInfo petsitter.SubmitSkillsRequest) error
	GetPetsitterStatus(userID uint) (*petsitter.PetSitterStatusResponse, error)
	FindPetSitterByID(id uint) (*entities.PetSitter, error)
	GetAllPetSitters(page, count int) (*petsitter.PetSittersListResponse, error)
	GetPetsitterServicesResponse(Services []enums.ServiceType) []entities.Service
	CheckPetSitterStatus(pss enums.PetSitterStatus) error
	CheckPetSitterStep(currentStep enums.OnboardingStep, requiredStep enums.OnboardingStep) error
	SubmitPersonalInfo(petSitterInfo petsitter.SubmitPersonalInfoRequest) error
	GetPetSitterDetails(info petsitter.GetPetSitterDetailsRequest) (*petsitter.PetSitterDetailsResponse, error)
	ChangePetSitterStatus(info petsitter.ChangePetSitterStatusRequest) error
}
