package service

import (
	"fmt"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/ports"
)

type RequestService struct {
	userService      usecase.UserService
	provinceService  usecase.ProvinceService
	cityService      usecase.CityService
	addressService   usecase.AddressService
	petService       usecase.PetService
	serviceService   usecase.ServiceService
	petSitterService usecase.PetSitterService
	unitOfWork       ports.UnitOfWork
}

func NewRequestService(userService usecase.UserService, unitOfWork ports.UnitOfWork, provinceService usecase.ProvinceService, addressService usecase.AddressService, petService usecase.PetService) *RequestService {
	return &RequestService{
		userService:     userService,
		unitOfWork:      unitOfWork,
		provinceService: provinceService,
		addressService:  addressService,
		petService:      petService,
	}
}

func (rs *RequestService) CreateRequest(info request.CreateRequestRequest) error {
	petSitterUser, err := rs.userService.FindUserByID(info.PetSitterUserID)
	if err != nil {
		return err
	}
	userRepo := rs.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(petSitterUser)
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

	err = rs.validateCalendarSlots(petSitter.Schedule, info.CalenderSlots)
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

	pets, err := rs.makeRequestPets(info.PetIDs)
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

	err = rs.validateRequestServices(petSitter, info.ServiceIDs)
	if err != nil {
		return err
	}

	services, err := rs.makeRequestServices(info.ServiceIDs)
	if err != nil {
		return err
	}

	totalPrice := rs.calculateTotalPrice(services, petSitter.Schedule, len(pets))

	calendarSlots := rs.makeCalendarSlots(info.CalenderSlots)

	newRequest := &entities.Request{
		UserID:          info.UserID,
		PetSitterUserID: info.PetSitterUserID,
		Status:          enums.Pending,
		Chat:            entities.Chat{},
		TransferID:      nil,
		CalendarSlots:   calendarSlots,
		Pets:            pets,
		TotalPrice:      uint(totalPrice),
		Notes:           info.Notes,
		Comment:         nil,
		Address:         *address,
		Services:        services,
	}

	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	err = requestRepo.CreateRequest(newRequest)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RequestService) GetCreateRequestInfo(info request.GetCreateRequestInfoRequest) (*request.CreateRequestInfoResponse, error) {
	addresses, err := rs.addressService.GetUserAddressesInfo(info.UserID)
	if err != nil {
		return nil, err
	}

	petsData, err := rs.petService.GetPetsBasicData(pet.GetPetsBasicDataRequest{
		UserID: info.UserID,
	})
	if err != nil {
		return nil, err
	}

	freeSlots, err := rs.petSitterService.GetPetSitterFreeSlotsResponse(info.PetSitterUserID)
	if err != nil {
		return nil, err
	}

	servicesData, err := rs.petSitterService.GetServicesResponse(info.PetSitterUserID)
	if err != nil {
		return nil, err
	}

	return &request.CreateRequestInfoResponse{
		Services:          servicesData,
		Addresses:         addresses,
		Pets:              petsData,
		FreeCalendarSlots: freeSlots,
	}, nil
}

func (rs *RequestService) validateCalendarSlots(petSitterSlots []entities.CalendarSlot, requestSlots []request.RequestCalendarSlotRequest) error {
	for _, req := range requestSlots {
		for _, userSlot := range req.Slots {
			flag := false
			for _, petSitterSlot := range petSitterSlots {
				if petSitterSlot.Date.Equal(req.Date) {
					for _, sitterSlot := range petSitterSlot.Slots {
						if sitterSlot == userSlot {
							flag = true
							break
						}
					}
				}
			}
			if !flag {
				return fmt.Errorf("invalid calendar slots")
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

func (rs *RequestService) makeRequestPets(petIDs []uint) ([]entities.Pet, error) {
	pets := make([]entities.Pet, 0)

	for _, petID := range petIDs {
		pet, err := rs.petService.FindPetByID(petID)
		if err != nil {
			return nil, err
		}
		requestPet := &entities.Pet{
			UserID:     pet.UserID,
			Name:       pet.Name,
			Kind:       pet.Kind,
			Species:    pet.Species,
			BirthDate:  pet.BirthDate,
			IsAdult:    pet.IsAdult,
			Gender:     pet.Gender,
			Weight:     pet.Weight,
			PictureKey: pet.PictureKey,
			AboutPet:   pet.AboutPet,
			Type:       "request",
		}

		pets = append(pets, *requestPet)
	}

	return pets, nil
}

func (rs *RequestService) validateRequestServices(petSitter *entities.PetSitter, serviceIDs []uint) error {
	for _, serviceID := range serviceIDs {
		flag := false
		for _, service := range petSitter.Services {
			if service.ID == serviceID {
				flag = true
				break
			}
		}
		if !flag {
			return fmt.Errorf("wrong service selected")
		}
	}

	return nil
}

func (rs *RequestService) makeRequestServices(serviceIDs []uint) ([]entities.Service, error) {
	services := make([]entities.Service, 0)

	for _, serviceID := range serviceIDs {
		service, err := rs.serviceService.FindServiceByID(serviceID)
		if err != nil {
			return nil, err
		}
		requestService := &entities.Service{
			PetSitterID: service.PetSitterID,
			Description: service.Description,
			Price:       service.Price,
			Type:        service.Type,
			PetKinds:    service.PetKinds,
			Kind:        "request",
		}

		services = append(services, *requestService)
	}

	return services, nil
}

func (rs *RequestService) calculateTotalPrice(services []entities.Service, slots []entities.CalendarSlot, petCount int) int {
	totalPrice := 0
	for _, service := range services {
		servicePrice := int(service.Price) * petCount
		totalPrice += servicePrice
	}

	for _, slot := range slots {
		slotHours := len(slot.Slots) / 2
		totalPrice *= slotHours
	}

	return totalPrice
}

func (rs *RequestService) makeCalendarSlots(calendarSlots []request.RequestCalendarSlotRequest) []entities.CalendarSlot {
	slots := make([]entities.CalendarSlot, 0)

	for _, slot := range calendarSlots {
		calendarSlot := &entities.CalendarSlot{
			Date:  slot.Date,
			Slots: slot.Slots,
		}
		slots = append(slots, *calendarSlot)
	}

	return slots
}
