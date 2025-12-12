package postgres

import (

	"gorm.io/gorm"
)

func applyModifiers(db *gorm.DB, modifiers ...QueryModifier) *gorm.DB {
	for _, m := range modifiers {
		db = m.Apply(db)
	}
	return db
}

// func applyQueryOptions(query *gorm.DB, options *dsl.ParsedQuery) *gorm.DB {
// 	if options == nil {
// 		return query
// 	}

// 	// ============== Filters ==============
// 	if len(options.Filters.Items) > 0 {
// 		for _, f := range options.Filters.Items {
// 			if f.Op == "IN" {
// 				query = query.Where(fmt.Sprintf("%s IN (?)", f.Field), f.Value)
// 				continue
// 			}

// 			query = query.Where(fmt.Sprintf("%s %s ?", f.Field, f.Op), f.Value)
// 		}
// 	}

// 	// ============== Sorting ==============
// 	if options.Sorts.Items != nil {
// 		for _, sort := range options.Sorts.Items {
// 			query = query.Order(sort.Field + " " + sort.Dir)
// 		}
// 		return query
// 	}

// 	// ============== Pagination ==============
// 	// if options.Pagination != nil {
// 	// 	query = query.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
// 	// }

// 	return query
// }

// type BaseRepository struct {
// 	db *gorm.DB
// }

// func NewBaseRepository(db *gorm.DB) *BaseRepository {
// 	return &BaseRepository{db: db}
// }

// func (r *BaseRepository) WithFilters(filters *dsl.Filters) *BaseRepository {
// 	// اعمال فیلترها به db
// 	db := dsl.ApplySQL(r.db, filters, nil)
// 	return &BaseRepository{db: db}
// }

// func (r *BaseRepository) WithSorts(sorts *dsl.Sorts) *BaseRepository {
// 	// اعمال سورت‌ها به db
// 	db := dsl.ApplySQL(r.db, nil , sorts)
// 	return &BaseRepository{db: db}
// }

// func (r *BaseRepository) WithPagination(pagination *dsl.Pagination) *BaseRepository {
// 	// اعمال صفحه‌بندی به db
// 	db := dsl.ApplySQL(r.db, pagination)
// 	return &BaseRepository{db: db}
// }

// func (r *BaseRepository) FindAll(out interface{}) error {
// 	if err := r.db.Find(out).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
