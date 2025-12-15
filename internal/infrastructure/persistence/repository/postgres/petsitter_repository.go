package postgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/dsl"

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

func (pr *PetSitterRepository) PreloadServices(petSitter *entities.PetSitter) error {
	return pr.db.Preload("Services").First(petSitter, petSitter.ID).Error
}

func (pr *PetSitterRepository) CreatePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Create(petSitter).Error
}

func (pr *PetSitterRepository) UpdatePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Save(petSitter).Error
}

func (pr *PetSitterRepository) FindPetSitterByUserID(id uint) (*entities.PetSitter, error) {
	var petsitter entities.PetSitter
	err := pr.db.First(&petsitter, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	// return &	, nil

	// if result := pr.db.First(&petsitter, "userid = ?", id); result.Error != nil {
	// 	if result.Error == gorm.ErrRecordNotFound {
	// 		return nil, nil
	// 	}
	// 	return nil, result.Error
	// }
	return &petsitter, nil
}

func (pr *PetSitterRepository) SearchPetSitters(offset int, limit int, filters []dsl.Filter, sorts []dsl.Sort) ([]*entities.PetSitter, int64, error) {
	var petSitters []*entities.PetSitter
	var total int64

	// Start with base query
	query := pr.db.Model(&entities.PetSitter{})

	// Apply filters
	for _, filter := range filters {
		query = applyFilterToQuery(query, filter)
	}

	// Count total records (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorts
	for _, sort := range sorts {
		query = query.Order(sort.Field + " " + sort.Dir)
	}

	// Apply pagination
	query = query.Offset(offset).Limit(limit)

	// Fetch paginated results
	if err := query.Find(&petSitters).Error; err != nil {
		return nil, 0, err
	}

	return petSitters, total, nil
}

// applyFilterToQuery applies a single filter to a gorm query
func applyFilterToQuery(query *gorm.DB, filter dsl.Filter) *gorm.DB {
	switch filter.Op {
	case "=":
		return query.Where(filter.Field+" = ?", filter.Value)
	case "!=":
		return query.Where(filter.Field+" != ?", filter.Value)
	case ">":
		return query.Where(filter.Field+" > ?", filter.Value)
	case "<":
		return query.Where(filter.Field+" < ?", filter.Value)
	case ">=":
		return query.Where(filter.Field+" >= ?", filter.Value)
	case "<=":
		return query.Where(filter.Field+" <= ?", filter.Value)
	case "LIKE":
		return query.Where(filter.Field+" LIKE ?", filter.Value)
	case "IN":
		return query.Where(filter.Field+" IN (?)", filter.Value)
	default:
		return query
	}
}
