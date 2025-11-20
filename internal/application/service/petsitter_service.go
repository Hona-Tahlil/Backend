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

func (ps *PetSitterService) GetPetSitterFreeSlotsResponse(id uint) ([]calendarslot.CalendarSlotInfoResponse, error) {
	user, err := ps.userService.FindUserByID(id, true)
	if err != nil {
		return nil, err
	}
	if user.PetSitter == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.PetSitter, bootstrap.Run().Constants.ErrorTags.NotFound)
		return nil, &ve
	}
	petSitter := user.PetSitter
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

func (ps *PetSitterService) GetServicesResponse(id uint) ([]servicedto.ServiceInfoResponse, error) {
	r := make([]servicedto.ServiceInfoResponse, 0)
	user, err := ps.userService.FindUserByID(id, true)
	if err != nil {
		return nil, err
	}
	petSitter := user.PetSitter
	if petSitter == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.PetSitter, bootstrap.Run().Constants.ErrorTags.NotFound)
		return nil, &ve
	}
	for _, service := range petSitter.Services {
		r = append(r, ps.serviceService.GetServiceResponse(&service))
	}
	return r, nil
}
