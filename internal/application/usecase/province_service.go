package usecase

import (
	"hona/backend/internal/application/dto/provincecity"
)

type ProvinceService interface {
	GetAllProvincesResponse() ([]provincecity.ProvinceResponse, error)
	GetCitiesByProvinceName(info provincecity.GetProvinceCitiesRequest) ([]provincecity.CityResponse, error)
}
