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
	"log"

	"github.com/samber/lo"
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
	if petSitter.Schedule == nil {
		return nil, nil
	}
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
	if petSitter.Services == nil {
		return r, nil
	}
	for _, service := range petSitter.Services {
		r = append(r, ps.serviceService.GetServiceResponse(&service))
	}
	return r, nil
}

func (ps *PetSitterService) GetAvailableServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error) {
	r := make([]servicedto.ServiceInfoResponse, 0)
	if petSitter.Services == nil {
		return r, nil
	}
	for _, service := range petSitter.Services {
		if service.Price != 0 {
			r = append(r, ps.serviceService.GetServiceResponse(&service))
		}
	}
	return r, nil
}

func (ps *PetSitterService) GetPetSitterByUserID(id uint) (*entities.PetSitter, error) {
	user, err := ps.userService.FindUserByID(id)
	if err != nil {
		return nil, err
	}
	ps.userService.PreloadFields(user, []string{"PetSitter"})
	if user.PetSitter == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	if user.PetSitter.Status != enums.PSS_Active {
		log.Printf("[DEBUG] not active:")
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	return user.PetSitter, nil
}

func (ps *PetSitterService) GetPetSitterByID(id uint) (*entities.PetSitter, error) {
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundPetSitter, err := petSitterRepo.FindPetSitterByID(id)
	if err != nil {
		return nil, err
	}
	if foundPetSitter == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	return foundPetSitter, nil
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
			var ve exceptions.ValidationErrors
			ve.AddError(bootstrap.Run().Constants.ErrorFields.Pet, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
			return &ve
		}
	}

	return nil
}

func (ps *PetSitterService) ValidateService(services []entities.Service, serviceID uint) (*entities.Service, error) {
	for _, service := range services {
		if service.ID == serviceID && service.Price != 0 {
			return &entities.Service{
				PetSitterID: service.PetSitterID,
				Type:        service.Type,
				Price:       service.Price,
				Description: service.Description,
				Kind:        bootstrap.Run().Constants.EntityConstants.Request,
			}, nil
		}
	}
	var ve exceptions.ValidationErrors
	ve.AddError(bootstrap.Run().Constants.ErrorFields.Service, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
	return nil, &ve
}

// can be accept or cancel
func (ps *PetSitterService) AutoUpdateSlots(petSitter *entities.PetSitter, calendarSlots []entities.CalendarSlot, accept bool) error {
	newSlots := make([]entities.CalendarSlot, 0)
	petSitterSlots := petSitter.Schedule
	for _, petSitterSlot := range petSitterSlots {
		for _, calendarSlot := range calendarSlots {
			if petSitterSlot.Date.Equal(calendarSlot.Date) {
				if accept && petSitterSlot.Status == enums.Free {
					petSitterSlot.Slots, _ = lo.Difference(petSitterSlot.Slots, calendarSlot.Slots)
				} else if !accept && petSitterSlot.Status == enums.Booked {
					petSitterSlot.Slots, _ = lo.Difference(petSitterSlot.Slots, calendarSlot.Slots)
				}
				newSlots = append(newSlots, ps.makeSlot(&calendarSlot, accept))
			}
		}
	}
	petSitter.Schedule = append(petSitter.Schedule, newSlots...)
	ps.removeEmptySitterSlots(petSitter)
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	return petSitterRepo.EditPetSitter(petSitter)
}

func (ps *PetSitterService) removeEmptySitterSlots(petSitter *entities.PetSitter) {
	for i, slot := range petSitter.Schedule {
		if len(slot.Slots) == 0 {
			petSitter.Schedule = append(petSitter.Schedule[:i], petSitter.Schedule[i+1:]...)
		}
	}
}

func (ps *PetSitterService) makeSlot(calendarSlot *entities.CalendarSlot, accept bool) entities.CalendarSlot {
	var slot entities.CalendarSlot
	slot.Date = calendarSlot.Date
	if accept {
		slot.Status = enums.Booked
	} else {
		slot.Status = enums.Free
	}
	slot.Slots = calendarSlot.Slots
	return slot
}
