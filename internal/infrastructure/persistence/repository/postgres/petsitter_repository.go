package postgres

import (
	"fmt"
	"hona/backend/internal/application/dto/general"
	"hona/backend/internal/domain/entities"

	"github.com/lib/pq"
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
	// Use Save but it should work correctly with the proper array type
	// If there are issues with associations, specify columns to update
	return pr.db.Save(petSitter).Error
}

func (pr *PetSitterRepository) ReplaceSchedule(petSitter *entities.PetSitter, schedule []entities.CalendarSlot) error {
	return pr.db.Model(petSitter).Association("Schedule").Replace(schedule)
}


// func (pr *PetSitterRepository) SearchPetSitters(options *QueryOptions) ([]*entities.PetSitter, int64, error) {
// 	var petSitters []*entities.PetSitter
// 	var total int64
// 	query := pr.db.Model(&entities.PetSitter{}).
// 		Select("DISTINCT pet_sitters.*").
// 		Joins("JOIN users ON users.id = pet_sitters.user_id").
// 		Joins("LEFT JOIN services ON services.pet_sitter_id = pet_sitters.id AND services.kind = ?", "petSitter").
// 		Joins("LEFT JOIN addresses ON addresses.refer = users.id AND addresses.type = ?", "User").
// 		Joins("LEFT JOIN calendar_slots ON calendar_slots.refer = pet_sitters.id")

// 	if options.Filters != nil {
// 		filterModifier := NewFilterModifier(options.Filters.Filters)
// 		query = filterModifier.Apply(query)
// 	}
// 	// if err := query.Count(&total).Error; err != nil {
// 	// 	return nil, 0, err
// 	// }
// 	if err := query.Session(&gorm.Session{}).
// 		Distinct("pet_sitters.id").
// 		Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}
// 	if options.Sorting != nil {
// 		sortModifier := NewSortModifier(options.Sorting.Sorts)
// 		query = sortModifier.Apply(query)
// 	}
// 	if options.Pagination != nil {
// 		paginationModifier := NewPaginationModifier(options.Pagination.Offset, options.Pagination.Limit)
// 		query = paginationModifier.Apply(query)
// 	}
// 	if err := query.Find(&petSitters).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	return petSitters, total, nil
// }

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

func (pr *PetSitterRepository) CreateService(service *entities.Service) error {
	return pr.db.Create(service).Error
}

func (pr *PetSitterRepository) UpdateService(service *entities.Service) error {
	return pr.db.Save(service).Error
}

func (pr *PetSitterRepository) DeleteService(service *entities.Service) error {
	return pr.db.Delete(service).Error
}


func (pr *PetSitterRepository) SearchPetSitters(options *QueryOptions) ([]*entities.PetSitter, int64, error) {
	var (
		petSitters []*entities.PetSitter
		total      int64
	)

	// 1) Parse known filters from []general.Filter
	fs := parsePetSitterSearchFilters(options)

	q := pr.db.Model(&entities.PetSitter{}).
		Joins("JOIN users ON users.id = pet_sitters.user_id").
		Joins("LEFT JOIN addresses ON addresses.refer = users.id AND addresses.type = ?", "User")

	// -----------------------
	// Filters
	// -----------------------

	if fs.City != nil && *fs.City != "" {
		q = q.Where("addresses.city = ?", *fs.City)
	}

	if len(fs.PetKindsAny) > 0 {
		q = q.Where("pet_sitters.pet_kinds && ?", pq.Int64Array(intsToInt64(fs.PetKindsAny)))
	}

	if fs.Date != nil && fs.StartSlot != nil && fs.EndSlot != nil {
		slots := makeSlotRange(*fs.StartSlot, *fs.EndSlot)
		q = q.Where(`
			EXISTS (
				SELECT 1
				FROM calendar_slots cs
				WHERE cs.pet_sitter_id = pet_sitters.id
				  AND DATE(cs.date) = DATE(?)
				  AND cs.slots @> ?
			)
		`, *fs.Date, pq.Int64Array(intsToInt64(slots)))
	} else if fs.Date != nil {
		q = q.Where(`
			EXISTS (
				SELECT 1
				FROM calendar_slots cs
				WHERE cs.pet_sitter_id = pet_sitters.id
				  AND DATE(cs.date) = DATE(?)
			)
		`, *fs.Date)
	}


	if fs.ServiceType != nil {
		q = q.Joins(
			"LEFT JOIN services s ON s.pet_sitter_id = pet_sitters.id AND s.kind = ? AND s.type = ?",
			"petSitter",
			*fs.ServiceType,
		)
		q = q.Where("s.id IS NOT NULL")
		if fs.MinPrice != nil {
			q = q.Where("s.price >= ?", *fs.MinPrice)
		}
		if fs.MaxPrice != nil {
			q = q.Where("s.price <= ?", *fs.MaxPrice)
		}
	} else if fs.MinPrice != nil || fs.MaxPrice != nil {
		where := `
		EXISTS (
			SELECT 1
			FROM services sx
			WHERE sx.pet_sitter_id = pet_sitters.id
			  AND sx.kind = ?
	`
		args := []any{"petSitter"}

		if fs.MinPrice != nil {
			where += " AND sx.price >= ?"
			args = append(args, *fs.MinPrice)
		}
		if fs.MaxPrice != nil {
			where += " AND sx.price <= ?"
			args = append(args, *fs.MaxPrice)
		}

		where += ")"
		q = q.Where(where, args...)
	}


	if err := q.Session(&gorm.Session{}).
		Distinct("pet_sitters.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if options != nil && options.Sorting != nil && len(options.Sorting.Sorts) > 0 {
		var err error
		q, err = applyPetSitterUISort(q, options.Sorting.Sorts, fs.ServiceType)
		if err != nil {
			return nil, 0, err
		}
	} else {
		q = q.Order("pet_sitters.created_at DESC")
	}


	if options != nil && options.Pagination != nil {
		q = q.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
	}

	// -----------------------
	// Fetch (Distinct to avoid duplicates)
	// -----------------------
	// To avoid "SELECT DISTINCT ON must match ORDER BY" error in PostgreSQL,
	// we use a subquery approach: get ordered distinct IDs, then fetch full records
	subquerySQL := pr.db.Model(&entities.PetSitter{}).
		Select("pet_sitters.id").
		Joins("JOIN users ON users.id = pet_sitters.user_id").
		Joins("LEFT JOIN addresses ON addresses.refer = users.id AND addresses.type = ?", "User")

	// Apply same filters to subquery
	if fs.City != nil && *fs.City != "" {
		subquerySQL = subquerySQL.Where("addresses.city = ?", *fs.City)
	}
	if len(fs.PetKindsAny) > 0 {
		subquerySQL = subquerySQL.Where("pet_sitters.pet_kinds && ?", pq.Int64Array(intsToInt64(fs.PetKindsAny)))
	}
	if fs.Date != nil && fs.StartSlot != nil && fs.EndSlot != nil {
		slots := makeSlotRange(*fs.StartSlot, *fs.EndSlot)
		subquerySQL = subquerySQL.Where(`
			EXISTS (
				SELECT 1
				FROM calendar_slots cs
				WHERE cs.pet_sitter_id = pet_sitters.id
				  AND DATE(cs.date) = DATE(?)
				  AND cs.slots @> ?
			)
		`, *fs.Date, pq.Int64Array(intsToInt64(slots)))
	} else if fs.Date != nil {
		subquerySQL = subquerySQL.Where(`
			EXISTS (
				SELECT 1
				FROM calendar_slots cs
				WHERE cs.pet_sitter_id = pet_sitters.id
				  AND DATE(cs.date) = DATE(?)
			)
		`, *fs.Date)
	}

	if fs.ServiceType != nil {
		subquerySQL = subquerySQL.Joins(
			"LEFT JOIN services s ON s.pet_sitter_id = pet_sitters.id AND s.kind = ? AND s.type = ?",
			"petSitter",
			*fs.ServiceType,
		)
		if fs.MinPrice != nil {
			subquerySQL = subquerySQL.Where("s.price >= ?", *fs.MinPrice)
		}
		if fs.MaxPrice != nil {
			subquerySQL = subquerySQL.Where("s.price <= ?", *fs.MaxPrice)
		}
	} else if fs.MinPrice != nil || fs.MaxPrice != nil {
		where := `
		EXISTS (
			SELECT 1
			FROM services sx
			WHERE sx.pet_sitter_id = pet_sitters.id
			  AND sx.kind = ?
	`
		args := []any{"petSitter"}

		if fs.MinPrice != nil {
			where += " AND sx.price >= ?"
			args = append(args, *fs.MinPrice)
		}
		if fs.MaxPrice != nil {
			where += " AND sx.price <= ?"
			args = append(args, *fs.MaxPrice)
		}

		where += ")"
		subquerySQL = subquerySQL.Where(where, args...)
	}

	// Apply sorting and pagination to the subquery
	if options != nil && options.Sorting != nil && len(options.Sorting.Sorts) > 0 {
		var err error
		subquerySQL, err = applyPetSitterUISort(subquerySQL, options.Sorting.Sorts, fs.ServiceType)
		if err != nil {
			return nil, 0, err
		}
	} else {
		subquerySQL = subquerySQL.Order("pet_sitters.created_at DESC")
	}

	if options != nil && options.Pagination != nil {
		subquerySQL = subquerySQL.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
	}

	// Fetch full records using the subquery for ordering and pagination
	// 1) Get ordered ids (subquery already has filters + sort + pagination)
	// 1) Get ordered distinct ids (subquery already has filters + sort + pagination)
	var ids []uint
	if err := subquerySQL.Pluck("pet_sitters.id", &ids).Error; err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		return []*entities.PetSitter{}, total, nil
	}

	// 2) Fetch full records and preserve exact order using CASE
	caseSQL := "CASE pet_sitters.id "
	caseArgs := make([]any, 0, len(ids))
	for i, id := range ids {
		caseSQL += "WHEN ? THEN " + fmt.Sprint(i) + " "
		caseArgs = append(caseArgs, id)
	}
	caseSQL += "END"

	if err := pr.db.Model(&entities.PetSitter{}).
		Where("pet_sitters.id IN ?", ids).
		Order(gorm.Expr(caseSQL, caseArgs...)).
		Find(&petSitters).Error; err != nil {
		return nil, 0, err
	}

	return petSitters, total, nil
}

func applyPetSitterUISort(q *gorm.DB, sorts []general.Sort, serviceType *int) (*gorm.DB, error) {
	field := sorts[0].Field

	switch field {
	case "max_price":
		if serviceType != nil {
			// same-type price (NULLS LAST)
			return q.Order("s.price DESC NULLS LAST"), nil
		}
		// no type => max over all services
		return q.Order(`
			(SELECT MAX(sx.price)
			 FROM services sx
			 WHERE sx.pet_sitter_id = pet_sitters.id
			   AND sx.kind = 'petSitter') DESC NULLS LAST
		`), nil

	case "min_price":
		if serviceType != nil {
			return q.Order("s.price ASC NULLS LAST"), nil
		}
		return q.Order(`
			(SELECT MIN(sx.price)
			 FROM services sx
			 WHERE sx.pet_sitter_id = pet_sitters.id
			   AND sx.kind = 'petSitter') ASC NULLS LAST
		`), nil

	case "max_rate":
		// avg rating (NULLS LAST)
		return q.Order(`
			(SELECT AVG(c.rating)::float8
			 FROM comments c
			 WHERE c.pet_sitter_id = pet_sitters.id) DESC NULLS LAST
		`), nil

	case "max_comments":
		return q.Order(`
			(SELECT COUNT(*)
			 FROM comments c
			 WHERE c.pet_sitter_id = pet_sitters.id) DESC NULLS LAST
		`), nil

	default:
		// fallback
		return q.Order("pet_sitters.created_at DESC"), nil
	}
}
