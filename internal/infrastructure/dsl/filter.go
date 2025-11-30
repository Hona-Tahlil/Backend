package dsl

type Filter struct {
	Field string
	Op    string // =, !=, >, <, >=, <=, LIKE, IN
	Value any
}

type Filters struct {
	Items []Filter
}

func NewFilters() *Filters {
	return &Filters{Items: make([]Filter, 0)}
}

// Equal
func (f *Filters) Eq(field string, value any) *Filters {
	f.Items = append(f.Items, Filter{Field: field, Op: "=", Value: value})
	return f
}

// Not Equal
func (f *Filters) Ne(field string, value any) *Filters {
	f.Items = append(f.Items, Filter{Field: field, Op: "!=", Value: value})
	return f
}

// Less / Greater
func (f *Filters) Lt(field string, value any) *Filters {
	f.Items = append(f.Items, Filter{Field: field, Op: "<", Value: value})
	return f
}

func (f *Filters) Gt(field string, value any) *Filters {
	f.Items = append(f.Items, Filter{Field: field, Op: ">", Value: value})
	return f
}

func (f *Filters) Lte(field string, value any) *Filters {
	f.Items = append(f.Items, Filter{Field: field, Op: "<=", Value: value})
	return f
}

func (f *Filters) Gte(field string, value any) *Filters {
	f.Items = append(f.Items, Filter{Field: field, Op: ">=", Value: value})
	return f
}

// LIKE
func (f *Filters) Like(field string, value string) *Filters {
	f.Items = append(f.Items, Filter{
		Field: field,
		Op:    "LIKE",
		Value: "%" + value + "%",
	})
	return f
}

// IN
func (f *Filters) In(field string, values []any) *Filters {
	f.Items = append(f.Items, Filter{Field: field, Op: "IN", Value: values})
	return f
}

// Applies filters to GORM
// func (f *Filters) Apply(db *gorm.DB) *gorm.DB {
// 	for _, flt := range f.items {

// 		if flt.Op == "IN" {
// 			db = db.Where(fmt.Sprintf("%s IN (?)", flt.Field), flt.Value)
// 			continue
// 		}

// 		db = db.Where(fmt.Sprintf("%s %s ?", flt.Field, flt.Op), flt.Value)
// 	}

// 	return db
// }

/////////////////////////////////
// type Server struct {
//     Host string
//     Port int
// }

// type Option func(*Server)

// func WithHost(host string) Option {
//     return func(s *Server) {
//         s.Host = host
//     }
// }

// func WithPort(port int) Option {
//     return func(s *Server) {
//         s.Port = port
//     }
// }

// func NewServer(opts ...Option) *Server {
//     s := &Server{
//         Host: "localhost", // مقدار پیش‌فرض
//         Port: 8080,
//     }

//     for _, opt := range opts {
//         opt(s)
//     }

//     return s
// }
