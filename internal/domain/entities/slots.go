package entities

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"hona/backend/internal/domain/enums"
)

type Slots []enums.Slot

func (s Slots) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	parts := make([]string, len(s))
	for i, slot := range s {
		parts[i] = strconv.FormatUint(uint64(slot), 10)
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ",")), nil
}

func (s *Slots) Scan(value interface{}) error {
	if value == nil {
		*s = Slots{}
		return nil
	}
	var raw string
	switch v := value.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("unsupported slots type: %T", value)
	}
	raw = strings.TrimSpace(raw)
	if raw == "{}" || raw == "" {
		*s = Slots{}
		return nil
	}
	raw = strings.TrimPrefix(raw, "{")
	raw = strings.TrimSuffix(raw, "}")
	if raw == "" {
		*s = Slots{}
		return nil
	}
	items := strings.Split(raw, ",")
	slots := make(Slots, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		val, err := strconv.ParseUint(item, 10, 64)
		if err != nil {
			return err
		}
		slots = append(slots, enums.Slot(val))
	}
	*s = slots
	return nil
}
