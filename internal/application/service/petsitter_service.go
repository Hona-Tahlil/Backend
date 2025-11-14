package service

import (
	"fmt"
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/ports"
)

type PetSitterService struct {
	unitOfWork  ports.UnitOfWork
	userService usecase.UserService
}

func NewPetSitterService(unitOfWork ports.UnitOfWork, userService usecase.UserService) *PetSitterService {
	return &PetSitterService{
		unitOfWork:  unitOfWork,
		userService: userService,
	}
}

func (ps *PetSitterService) GetPetSitterFreeSlotsResponse(id uint) ([]calendarslot.CalendarSlotInfoResponse, error) {
	user, err := ps.userService.FindUserByID(id)
	if err != nil {
		return nil, err
	}
	userRepo := ps.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(user)
	if err != nil {
		return nil, err
	}
	if user.PetSitter == nil {
		return nil, fmt.Errorf("wrong pet sitter")
	}
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	err = userRepo.PreloadPetSitter(user)
	if err != nil {
		return nil, err
	}
	petSitter := user.PetSitter
	err = petSitterRepo.PreloadSchedule(petSitter)
	if err != nil {
		return nil, err
	}
	petSitterCalendarSlots := petSitter.Schedule
	freeSlots := make([]calendarslot.CalendarSlotInfoResponse, 0)
	for _, slot := range petSitterCalendarSlots {
		if slot.Status != enums.Free {
			continue
		}
		freeSlots = append(freeSlots, calendarslot.CalendarSlotInfoResponse{
			ID:    slot.ID,
			Date:  slot.Date,
			Slots: slot.Slots,
		})
	}
	return freeSlots, nil
}

func (ps *PetSitterService) GetServicesResponse(id uint) ([]servicedto.ServiceInfoResponse, error) {
	r := make([]servicedto.ServiceInfoResponse, 0)
	userRepo := ps.unitOfWork.Factory().UserRepository()
	user, err := userRepo.FindUserByID(id)
	if err != nil {
		return nil, err
	}
	err = userRepo.PreloadPetSitter(user)
	if err != nil {
		return nil, err
	}
	petSitter := user.PetSitter
	if petSitter == nil {
		return nil, fmt.Errorf("invalid pet sitter")
	}
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	err = petSitterRepo.PreloadServices(petSitter)
	if err != nil {
		return nil, err
	}
	for _, service := range petSitter.Services {
		r = append(r, servicedto.ServiceInfoResponse{
			ID:          service.ID,
			Type:        service.Type.String(),
			Description: service.Description,
			Price:       service.Price,
		})
	}
	return r, nil
}
