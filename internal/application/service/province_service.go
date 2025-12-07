package service

import (
	"hona/backend/internal/application/dto/provincecity"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/ports"
)

type ProvinceService struct {
	unitOfWork ports.UnitOfWork
}

func NewProvinceService(unitOfWork ports.UnitOfWork) *ProvinceService {
	return &ProvinceService{
		unitOfWork: unitOfWork,
	}
}

func (ps *ProvinceService) GetAllProvincesResponse() ([]provincecity.ProvinceResponse, error) {
	provinces := enums.GetAllProvinces()

	provinceResponses := make([]provincecity.ProvinceResponse, 0)
	for _, province := range provinces {
		provinceResponses = append(provinceResponses, provincecity.ProvinceResponse{
			Num:  province,
			Name: province.String(),
		})
	}

	return provinceResponses, nil
}

func (ps *ProvinceService) GetCitiesByProvinceName(info provincecity.GetProvinceCitiesRequest) ([]provincecity.CityResponse, error) {
	cities := enums.ProvinceWithCities[info.ProvinceNum]
	cityResponses := make([]provincecity.CityResponse, 0)
	for _, city := range cities {
		cityResponses = append(cityResponses, provincecity.CityResponse{
			Num:  city,
			Name: city.String(),
		})
	}

	return cityResponses, nil
}
