package service

import (
	"hona/backend/bootstrap"
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
)

type PetSitterService struct {
	unitOfWork          ports.UnitOfWork
	userService         usecase.UserService
	serviceService      usecase.ServiceService
	calendarSlotService usecase.CalendarSlotService
}

func NewPetSitterService(unitOfWork ports.UnitOfWork, userService usecase.UserService, serviceService usecase.ServiceService, calendarSlotService usecase.CalendarSlotService) *PetSitterService {
	return &PetSitterService{
		unitOfWork:          unitOfWork,
		userService:         userService,
		serviceService:      serviceService,
		calendarSlotService: calendarSlotService,
	}
}

func (ps *PetSitterService) GetPetSitterFreeSlotsResponse(petSitter *entities.PetSitter) ([]calendarslot.CalendarSlotInfoResponse, error) {
	petSitterCalendarSlots := petSitter.Schedule
	freeSlots := make([]entities.CalendarSlot, 0)
	for _, slot := range petSitterCalendarSlots {
		if slot.Status != enums.Free {
			continue
		}
		freeSlots = append(freeSlots, slot)
	}

	return ps.calendarSlotService.GetCalendarSlotsResponse(freeSlots), nil
}

func (ps *PetSitterService) GetServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error) {
	r := make([]servicedto.ServiceInfoResponse, 0)
	for _, service := range petSitter.Services {
		r = append(r, ps.serviceService.GetServiceResponse(&service))
	}
	return r, nil
}

func (ps *PetSitterService) GetAvailableServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error) {
	r := make([]servicedto.ServiceInfoResponse, 0)
	for _, service := range petSitter.Services {
		if service.Price != 0 {
			r = append(r, ps.serviceService.GetServiceResponse(&service))
		}
	}
	return r, nil
}

func (ps *PetSitterService) GetPetSitterByID(id uint) (*entities.PetSitter, error) {
	user, err := ps.userService.FindUserByID(id)
	if err != nil {
		return nil, err
	}
	ps.userService.PreloadFields(user, []string{"PetSitter"})
	if user.PetSitter == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	if !user.PetSitter.IsVerified {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	return user.PetSitter, nil
}

func (ps *PetSitterService) PreloadFields(petSitter *entities.PetSitter, fields []string) error {
	return ps.unitOfWork.Factory().PetSitterRepository().PreloadFields(petSitter, fields)
}

func (ps *PetSitterService) ValidatePets(pets []entities.Pet, petKinds []enums.PetKind) error {
	for _, pet := range pets {
		flag := false
		for _, kind := range petKinds {
			if kind == pet.Kind {
				flag = true
				break
			}
		}
		if !flag {
			var ce exceptions.ConflictErrors
			ce.Add(bootstrap.Run().Constants.ErrorFields.Pet, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
			return &ce
		}
	}

	return nil
}

func (ps *PetSitterService) ValidateService(services []entities.Service, serviceID uint) (*entities.Service, error) {
	for _, service := range services {
		if service.ID == serviceID && service.Price != 0 {
			return &service, nil
		}
	}
	var ce exceptions.ConflictErrors
	ce.Add(bootstrap.Run().Constants.ErrorFields.Service, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
	return nil, &ce
}
