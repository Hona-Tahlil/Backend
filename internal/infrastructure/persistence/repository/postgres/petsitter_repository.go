package postgres

import (
	"hona/backend/internal/domain/entities"
	domainpostgres "hona/backend/internal/domain/ports/postgres"

	"gorm.io/gorm"
)

type PetSitterRepository struct {
	db *gorm.DB
}

func NewPetSitterRepository(db *gorm.DB) *PetSitterRepository {
	return &PetSitterRepository{
		db: db,
	}
}

func (pr *PetSitterRepository) FindPetSitterByUserID(id uint) (*entities.PetSitter, error) {
	var foundPetsitter entities.PetSitter

	// if result := pr.db.First(&foundPetsitter, "userID = ?", id); result.Error != nil {
	if result := pr.db.Where("user_id = ?", id).First(&foundPetsitter); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundPetsitter, nil
}

func (pr *PetSitterRepository) PreloadFields(petSitter *entities.PetSitter, fields []string) error {
	for _, field := range fields {
		err := pr.db.Preload(field).First(petSitter, petSitter.ID).Error
		if err != nil {
			return err
		}
	}

	return nil
}
func (pr *PetSitterRepository) PreloadServices(petSitter *entities.PetSitter) error {
	return pr.db.Preload("Services").First(petSitter, petSitter.ID).Error
}

func (pr *PetSitterRepository) CreatePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Create(petSitter).Error
}

func (pr *PetSitterRepository) UpdatePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Save(petSitter).Error
}

func (pr *PetSitterRepository) SearchPetSitters(options *domainpostgres.QueryOptions) ([]*entities.PetSitter, int64, error) {
	var petSitters []*entities.PetSitter
	var total int64
	query := pr.db.Model(&entities.PetSitter{})
	if options.Filters != nil {
		filterModifier := NewFilterModifier(options.Filters.Filters)
		query = filterModifier.Apply(query)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if options.Sorting != nil {
		sortModifier := NewSortModifier(options.Sorting.Sorts)
		query = sortModifier.Apply(query)
	}
	if options.Pagination != nil {
		paginationModifier := NewPaginationModifier(options.Pagination.Offset, options.Pagination.Limit)
		query = paginationModifier.Apply(query)
	}
	if err := query.Find(&petSitters).Error; err != nil {
		return nil, 0, err
	}

	return petSitters, total, nil
}

func (pr *PetSitterRepository) FindPetSitterByID(id uint) (*entities.PetSitter, error) {
	var petSitter entities.PetSitter
	result := pr.db.First(&petSitter, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &petSitter, nil
}

func (pr *PetSitterRepository) GetAllPetSitters(limit, offset int) ([]entities.PetSitter, error) {
	var petSitters []entities.PetSitter
	err := pr.db.Limit(limit).Offset(offset).Find(&petSitters).Error
	if err != nil {
		return nil, err
	}
	return petSitters, nil
}

func (pr *PetSitterRepository) GetPetSittersCount() (int64, error) {
	var count int64
	err := pr.db.Model(&entities.PetSitter{}).Count(&count).Error
	return count, err
}

func (pr *PetSitterRepository) FindServiceByID(id uint) (*entities.Service, error) {
	var foundService entities.Service

	if result := pr.db.First(&foundService, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundService, nil
}
