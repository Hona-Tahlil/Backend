package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
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

func (ps *ProvinceService) FindProvinceByName(name enums.Province) (*entities.Province, error) {
	provinceRepo := ps.unitOfWork.Factory().ProvinceRepository()
	province, err := provinceRepo.FindProvinceByName(name)
	if err != nil {
		return nil, err
	}
	if province == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Province, bootstrap.Run().Constants.ErrorTags.NotFound)
		return nil, &ve
	}

	return province, nil
}
