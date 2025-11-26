package usecase

import (
	"hona/backend/internal/application/dto/provincecity"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
)

type ProvinceService interface {
	FindProvinceByName(name enums.Province) (*entities.Province, error)
	GetAllProvincesResponse() ([]provincecity.ProvinceResponse, error)
	GetCitiesByProvinceName(info provincecity.GetProvinceCitiesRequest) ([]provincecity.CityResponse, error)
}
