package service

import (
	"fmt"
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
)

type RequestService struct {
	userService usecase.UserService
}

func NewRequestService(userService usecase.UserService) *RequestService {
	return &RequestService{
		userService: userService,
	}
}
func (rs *RequestService) CreateRequest(info request.CreateRequestRequest) error {
	petSitterUser, err := rs.userService.FindUserByID(info.PetSitterUserID)
	if err != nil {
		return err
	}
	if petSitterUser.PetSitter == nil {
		return fmt.Errorf("invalid pet sitter")
	}
	petSitter := petSitterUser.PetSitter
	if !petSitter.IsVerified {
		return fmt.Errorf("invalid pet sitter")
	}
	// TODO: preload Schedule
	// err = rs.ValidateCalendarSlots(petSitter.Schedule, info.CalenderSlots)
	// if err != nil {
	// 	return err
	// }

	user, err := rs.userService.FindUserByID(info.UserID)
	if err != nil {
		return err
	}
	if !user.IsEmailVerified {
		return fmt.Errorf("invalid user")
	}
	pets := make([]entities.Pet, 0)
	// TODO: Preload Pets
	for _, petID := range info.PetIDs {
		flag := false
		for _, pet := range user.Pets {
			if pet.ID == petID {
				flag = true
				pets = append(pets, pet)
				break
			}
		}
		if !flag {
			return fmt.Errorf("wrong pet selected")
		}
	}

	for _, pet := range pets {
		flag := false
		for _, kind := range petSitter.PetKinds {
			if kind == pet.Kind {
				flag = true
				break
			}
		}
		if !flag {
			return fmt.Errorf("wrong pet selected")
		}
	}

	// TODO: find request Address by ID if not found search Addresses

	return nil
}

// func (rs *RequestService) ValidateCalendarSlots(petSitterSlots []entities.CalendarSlot, requestSlots []request.RequestCalendarSlotRequest) error {
// 	sort.Slice(requestSlots, func(i, j int) bool {
// 		return requestSlots[i].StartTime.Before(requestSlots[j].StartTime)
// 	})
// 	for i, slot := range requestSlots {
// 		if slot.EndTime.After(slot.StartTime) {
// 			return fmt.Errorf("slot start and end invalid")
// 		}
// 		prev := requestSlots[i-1]
// 		if slot.StartTime.Before(prev.EndTime) {
// 			return fmt.Errorf("slots overlap")
// 		}
// 	}

// 	for _, slot := range requestSlots {
// 		flag := false
// 		for _, sitterSlot := range petSitterSlots {
// 			if sitterSlot.Status != enums.Free {
// 				continue
// 			}
// 			if slot.StartTime.After(sitterSlot.StartTime) && slot.EndTime.Before(sitterSlot.EndTime) {
// 				flag = true
// 				break
// 			}
// 		}
// 		if !flag {
// 			// TODO: custom conflict error
// 			return fmt.Errorf("Pet Sitter is not free at that time")
// 		}
// 	}

// 	return nil
// }
