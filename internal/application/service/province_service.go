package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/provincecity"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	"log"
)

type ProvinceService struct {
	unitOfWork ports.UnitOfWork
}

func NewProvinceService(unitOfWork ports.UnitOfWork) *ProvinceService {
	return &ProvinceService{
		unitOfWork: unitOfWork,
	}
}

func (ps *ProvinceService) FindProvinceByName(name enums.Province) (*entities.Province, error) {
	provinceRepo := ps.unitOfWork.Factory().ProvinceRepository()
	province, err := provinceRepo.FindProvinceByName(name)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] province: %v", province)

	if province == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Province, bootstrap.Run().Constants.ErrorTags.NotFound)
		return nil, &ve
	}

	return province, nil
}
func (ps *ProvinceService) GetAllProvincesResponse() ([]provincecity.ProvinceResponse, error) {
	provinceRepo := ps.unitOfWork.Factory().ProvinceRepository()
	provinces, err := provinceRepo.GetAllProvinces()
	if err != nil {
		return nil, err
	}

	provinceResponses := make([]provincecity.ProvinceResponse, 0)
	for _, province := range provinces {
		provinceResponses = append(provinceResponses, provincecity.ProvinceResponse{
			Num:  province.Name,
			Name: province.Name.String(),
		})
	}

	return provinceResponses, nil
}

func (ps *ProvinceService) GetCitiesByProvinceName(info provincecity.GetProvinceCitiesRequest) ([]provincecity.CityResponse, error) {
	province, err := ps.FindProvinceByName(info.ProvinceNum)
	if err != nil {
		return nil, err
	}
	cityResponses := make([]provincecity.CityResponse, 0)
	for _, city := range province.Cities {
		cityResponses = append(cityResponses, provincecity.CityResponse{
			Num:  city.Name,
			Name: city.Name.String(),
		})
	}

	return cityResponses, nil
}
