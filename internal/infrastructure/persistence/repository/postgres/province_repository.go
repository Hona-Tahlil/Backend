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

	if result := pr.db.Preload("Cities").First(&province, "name = ?", name); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &province, nil
}

func (pr *ProvinceRepository) CreateProvince(province *entities.Province) error {
	return pr.db.Create(province).Error
}

func (pr *ProvinceRepository) GetAllProvinces() ([]entities.Province, error) {
	var provinces []entities.Province
	if err := pr.db.Preload("Cities").Find(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}
