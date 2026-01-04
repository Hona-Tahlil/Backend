package postgres

import (
	"hona/backend/internal/application/dto/general"
	"time"

	"gorm.io/gorm"
)

type FilterModifier struct {
	Filters []general.Filter
}

func NewFilterModifier(filters []general.Filter) FilterModifier {
	return FilterModifier{Filters: filters}
}

func (f FilterModifier) Apply(db *gorm.DB) *gorm.DB {
	for _, filter := range f.Filters {
		if filter.Op == "IN" {
			db = db.Where(filter.Field+" IN (?)", filter.Value)
			continue
		}
		db = db.Where(filter.Field+" "+filter.Op+" ?", filter.Value)
	}
	return db
}

type petSitterSearchFilters struct {
	City        *string
	ServiceType *int
	MinPrice    *uint
	MaxPrice    *uint

	PetKindsAny []int

	Date      *time.Time
	StartSlot *int
	EndSlot   *int
}

func parsePetSitterSearchFilters(options *QueryOptions) petSitterSearchFilters {
	var out petSitterSearchFilters
	if options == nil || options.Filters == nil {
		return out
	}

	for _, f := range options.Filters.Filters {
		switch f.Field {
		case "city":
			if v, ok := toString(f.Value); ok {
				out.City = &v
			}
		case "serviceType":
			if v, ok := toInt(f.Value); ok {
				out.ServiceType = &v
			}
		case "minPrice":
			if v, ok := toUint(f.Value); ok {
				out.MinPrice = &v
			}
		case "maxPrice":
			if v, ok := toUint(f.Value); ok {
				out.MaxPrice = &v
			}
		case "petKinds":
			if v, ok := toIntSlice(f.Value); ok {
				out.PetKindsAny = v
			}
		case "date":
			if t, ok := toTime(f.Value); ok {
				out.Date = &t
			}
		case "startSlot":
			if v, ok := toInt(f.Value); ok {
				out.StartSlot = &v
			}
		case "endSlot":
			if v, ok := toInt(f.Value); ok {
				out.EndSlot = &v
			}
		}
	}
	return out
}

func toString(v any) (string, bool) {
	s, ok := v.(string)
	return s, ok
}

func toInt(v any) (int, bool) {
	switch x := v.(type) {
	case int:
		return x, true
	case int32:
		return int(x), true
	case int64:
		return int(x), true
	case float64: // common from JSON decode
		return int(x), true
	default:
		return 0, false
	}
}

func toUint(v any) (uint, bool) {
	switch x := v.(type) {
	case uint:
		return x, true
	case uint32:
		return uint(x), true
	case uint64:
		return uint(x), true
	case int:
		if x < 0 {
			return 0, false
		}
		return uint(x), true
	case float64:
		if x < 0 {
			return 0, false
		}
		return uint(x), true
	default:
		return 0, false
	}
}

func toIntSlice(v any) ([]int, bool) {
	switch x := v.(type) {
	case []int:
		return x, true
	case []any:
		out := make([]int, 0, len(x))
		for _, it := range x {
			n, ok := toInt(it)
			if !ok {
				return nil, false
			}
			out = append(out, n)
		}
		return out, true
	case []float64:
		out := make([]int, 0, len(x))
		for _, it := range x {
			out = append(out, int(it))
		}
		return out, true
	default:
		return nil, false
	}
}

func toTime(v any) (time.Time, bool) {
	switch x := v.(type) {
	case time.Time:
		return x, true
	case string:
		t, err := time.Parse(time.RFC3339, x)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	default:
		return time.Time{}, false
	}
}

func makeSlotRange(start, end int) []int {
	if end < start {
		start, end = end, start
	}
	out := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		out = append(out, i)
	}
	return out
}

func intsToInt64(in []int) []int64 {
	out := make([]int64, 0, len(in))
	for _, x := range in {
		out = append(out, int64(x))
	}
	return out
}
