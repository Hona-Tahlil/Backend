package entities

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"hona/backend/internal/domain/enums"
)

type PetKinds []enums.PetKind

func (pk PetKinds) Value() (driver.Value, error) {
	if len(pk) == 0 {
		return "{}", nil
	}
	parts := make([]string, len(pk))
	for i, kind := range pk {
		parts[i] = strconv.FormatUint(uint64(kind), 10)
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ",")), nil
}

func (pk *PetKinds) Scan(value interface{}) error {
	if value == nil {
		*pk = PetKinds{}
		return nil
	}
	var raw string
	switch v := value.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("unsupported pet kinds type: %T", value)
	}
	raw = strings.TrimSpace(raw)
	if raw == "{}" || raw == "" {
		*pk = PetKinds{}
		return nil
	}
	raw = strings.TrimPrefix(raw, "{")
	raw = strings.TrimSuffix(raw, "}")
	if raw == "" {
		*pk = PetKinds{}
		return nil
	}
	items := strings.Split(raw, ",")
	kinds := make(PetKinds, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		val, err := strconv.ParseUint(item, 10, 64)
		if err != nil {
			return err
		}
		kinds = append(kinds, enums.PetKind(val))
	}
	*pk = kinds
	return nil
}
