package service

import (
	"fmt"
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/ports"
)

type RequestService struct {
	userService     usecase.UserService
	provinceService usecase.ProvinceService
	cityService     usecase.CityService
	addressService  usecase.AddressService
	unitOfWork      ports.UnitOfWork
}

func NewRequestService(userService usecase.UserService, unitOfWork ports.UnitOfWork, provinceService usecase.ProvinceService, addressService usecase.AddressService) *RequestService {
	return &RequestService{
		userService:     userService,
		unitOfWork:      unitOfWork,
		provinceService: provinceService,
		addressService:  addressService,
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

	petSitterRepo := rs.unitOfWork.Factory().PetSitterRepository()
	err = petSitterRepo.PreloadSchedule(petSitter)
	if err != nil {
		return err
	}

	err = rs.ValidateCalendarSlots(petSitter.Schedule, info.CalenderSlots)
	if err != nil {
		return err
	}

	user, err := rs.userService.FindUserByID(info.UserID)
	if err != nil {
		return err
	}
	if !user.IsEmailVerified {
		return fmt.Errorf("invalid user")
	}

	err = rs.validateRequestPets(user, petSitter, info.PetIDs)
	if err != nil {
		return err
	}

	var address *entities.Address
	if info.AddressInfo != nil {
		province, err := rs.provinceService.FindProvinceByName(info.AddressInfo.ProvinceName)
		if err != nil {
			return err
		}

		city, err := rs.cityService.FindCityByNameInProvince(info.AddressInfo.CityName, province)
		if err != nil {
			return err
		}

		address = &entities.Address{
			Province:      *province,
			City:          *city,
			StreetAddress: info.AddressInfo.StreetAddress,
			HouseNumber:   info.AddressInfo.HouseNumber,
			Unit:          info.AddressInfo.Unit,
			PostalCode:    info.AddressInfo.PostalCode,
		}
	} else if info.AddressID != nil {
		address, err = rs.addressService.FindAddressByID(*info.AddressID)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("either address info or address ID must be provided")
	}

	_ = address

	// requestRepo := rs.unitOfWork.Factory().RequestRepository()
	// err = requestRepo.CreateRequest(info, user, petSitterUser, address)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (rs *RequestService) ValidateCalendarSlots(petSitterSlots []entities.CalendarSlot, requestSlots []request.RequestCalendarSlotRequest) error {
	freeSlots := make(map[enums.Slot]bool)

	for _, sitter := range petSitterSlots {
		if sitter.Status != enums.Free {
			continue
		}
		for _, slot := range sitter.Slots {
			freeSlots[slot] = true
		}
	}

	for _, req := range requestSlots {
		for _, userSlot := range req.Slots {
			if ok, found := freeSlots[userSlot]; !found || !ok {
				return fmt.Errorf("pet sitter is not free at requested time: %v", userSlot)
			}
		}
	}

	return nil
}

func (rs *RequestService) validateRequestPets(user *entities.User, petSitter *entities.PetSitter, petIDs []uint) error {
	pets := make([]entities.Pet, 0)
	userRepo := rs.unitOfWork.Factory().UserRepository()
	err := userRepo.PreloadPets(user)
	if err != nil {
		return err
	}

	for _, petID := range petIDs {
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

	return nil
}
