package postgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type ProvinceRepository struct {
	db *gorm.DB
}

func NewProvinceRepository(db *gorm.DB) *ProvinceRepository {
	return &ProvinceRepository{
		db: db,
	}
}

func (pr *ProvinceRepository) FindProvinceByName(name enums.Province) (*entities.Province, error) {
	var province entities.Province

	if result := pr.db.First(&province, "name = ?", name); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	err := pr.db.Preload("Cities").First(province, province.ID).Error
	if err != nil {
		return nil, err
	}

	return &province, nil
}
