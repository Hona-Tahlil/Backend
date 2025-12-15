package dsl


type Sort struct {
	Field string
	Dir   string // ASC or DESC
}

type Sorts struct {
	Items []Sort
}

func NewSorts() *Sorts {
	return &Sorts{Items: make([]Sort, 0)}
}

func (s *Sorts) Asc(field string) *Sorts {
	s.Items = append(s.Items, Sort{Field: field, Dir: "ASC"})
	return s
}

func (s *Sorts) Desc(field string) *Sorts {
	s.Items = append(s.Items, Sort{Field: field, Dir: "DESC"})
	return s
}

// func (s *Sorts) Apply(db *gorm.DB) *gorm.DB {
// 	for _, sort := range s.Items {
// 		db = db.Order(sort.Field + " " + sort.Dir)
// 	}
// 	return db
// }
