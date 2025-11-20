package dsl

import (
	"fmt"
	"gorm.io/gorm"
)

type Filter struct {
	Field string
	Op    string   // =, !=, >, <, >=, <=, LIKE, IN
	Value any
}

type Filters struct {
	items []Filter
}

func NewFilters() *Filters {
	return &Filters{items: make([]Filter, 0)}
}

// Equal
func (f *Filters) Eq(field string, value any) *Filters {
	f.items = append(f.items, Filter{Field: field, Op: "=", Value: value})
	return f
}

// Not Equal
func (f *Filters) Ne(field string, value any) *Filters {
	f.items = append(f.items, Filter{Field: field, Op: "!=", Value: value})
	return f
}

// Less / Greater
func (f *Filters) Lt(field string, value any) *Filters {
	f.items = append(f.items, Filter{Field: field, Op: "<", Value: value})
	return f
}

func (f *Filters) Gt(field string, value any) *Filters {
	f.items = append(f.items, Filter{Field: field, Op: ">", Value: value})
	return f
}

func (f *Filters) Lte(field string, value any) *Filters {
	f.items = append(f.items, Filter{Field: field, Op: "<=", Value: value})
	return f
}

func (f *Filters) Gte(field string, value any) *Filters {
	f.items = append(f.items, Filter{Field: field, Op: ">=", Value: value})
	return f
}

// LIKE
func (f *Filters) Like(field string, value string) *Filters {
	f.items = append(f.items, Filter{
		Field: field,
		Op:    "LIKE",
		Value: "%" + value + "%",
	})
	return f
}

// IN
func (f *Filters) In(field string, values []any) *Filters {
	f.items = append(f.items, Filter{Field: field, Op: "IN", Value: values})
	return f
}

// Applies filters to GORM
func (f *Filters) Apply(db *gorm.DB) *gorm.DB {
	for _, flt := range f.items {

		if flt.Op == "IN" {
			db = db.Where(fmt.Sprintf("%s IN (?)", flt.Field), flt.Value)
			continue
		}

		db = db.Where(fmt.Sprintf("%s %s ?", flt.Field, flt.Op), flt.Value)
	}

	return db
}
