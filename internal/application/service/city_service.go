package service

import (
	"fmt"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/ports"
)

type CityService struct {
	unitOfWork ports.UnitOfWork
}

func NewCityService(unitOfWork ports.UnitOfWork) *CityService {
	return &CityService{
		unitOfWork: unitOfWork,
	}
}

func (cs *CityService) FindCityByNameInProvince(name enums.City, province *entities.Province) (*entities.City, error) {
	cityRepo := cs.unitOfWork.Factory().CityRepository()
	city, err := cityRepo.FindCityByName(name)
	if err != nil {
		return nil, err
	}
	if city == nil {
		return nil, fmt.Errorf("invalid city")
	}

	flag := false
	for _, c := range province.Cities {
		if c.ID == city.ID {
			flag = true
			break
		}
	}
	if !flag {
		return nil, fmt.Errorf("city does not belong to province")
	}

	return city, nil
}
