package usecase

import (
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
)

type PetSitterService interface {
	GetCalendarSlots(info calendarslot.GetCalendarSlotsRequest) ([]calendarslot.CalendarSlotInfoResponse, error)
	UpdateFreeCalendarSlots(info calendarslot.UpdateFreeCalendarSlotsRequest) error
	GetPetSitterFreeSlotsResponse(petSitter *entities.PetSitter) ([]calendarslot.CalendarSlotInfoResponse, error)
	GetServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error)
	GetAvailableServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error)
	GetPetSitterByID(id uint) (*entities.PetSitter, error)
	UpdateRatingAndCommentsCount(petSitterID uint, rating float32, commentsCount uint) error
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
	GetPetKinds(info petsitter.GetPetKindsRequest) ([]pet.PetKindResponse, error)
	UpdatePetKinds(info petsitter.UpdatePetKindsRequest) error
	GetServices(info petsitter.GetServicesRequest) ([]servicedto.ServiceInfoResponse, error)
	CreateService(info petsitter.CreateServiceRequest) (servicedto.ServiceInfoResponse, error)
	UpdateService(info petsitter.UpdateServiceRequest) (servicedto.ServiceInfoResponse, error)
	DeleteService(info petsitter.DeleteServiceRequest) error
	FindPetSitterByID(id uint) (*entities.PetSitter, error)
	GetAllPetSitters(page, count int) (*petsitter.PetSittersListResponse, error)
	GetPetsitterServicesResponse(Services []enums.ServiceType, petSitterID uint) []entities.Service
	CheckPetSitterStatus(pss enums.PetSitterStatus) error
	CheckPetSitterStep(currentStep enums.OnboardingStep, requiredStep enums.OnboardingStep) error
	SubmitPersonalInfo(petSitterInfo petsitter.SubmitPersonalInfoRequest) error
	GetCalendarSlotsResponse(calendarSlots []entities.CalendarSlot) []calendarslot.CalendarSlotInfoResponse
	GetFreeMap(calendarSlots []entities.CalendarSlot) map[string]map[interface{}]bool
	FindServiceByID(id uint) (*entities.Service, error)
	GetServiceResponse(serviceEntity *entities.Service) servicedto.ServiceInfoResponse
	SearchPetSitters(info petsitter.SearchPetSittersRequest) ([]*petsitter.PetSitterInfoResponse, int64, error)
	SearchPetSittersForAdmin(info petsitter.AdminSearchPetSittersRequest) ([]petsitter.PetSitterListItemResponse, int64, error)
	GetPetSitterDetails(info petsitter.GetPetSitterDetailsRequest) (*petsitter.PetSitterDetailsResponse, error)
	GetPetSitterProfile(info petsitter.GetPetSitterProfileRequest) (*petsitter.PetSitterProfileResponse, error)
	ChangePetSitterStatus(info petsitter.ChangePetSitterStatusRequest) error
}
