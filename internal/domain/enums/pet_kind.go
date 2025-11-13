package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type PetKind uint

const (
	Dog PetKind = iota + 1
	Cat
	Bird
)

func (petKind PetKind) String() string {
	switch petKind {
	case PetKind(Dog):
		return "سگ"
	case PetKind(Cat):
		return "گربه"
	case PetKind(Bird):
		return "پرنده"
	}
	return ""
}

func GetAllPetKinds() []PetKind {
	return []PetKind{
		Dog,
		Cat,
		Bird,
	}
}

type PetKindSlice []PetKind

func (PetKindSlice) GormDataType() string {
	return "jsonb"
}

func (p PetKindSlice) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PetKindSlice) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %T to []byte", value)
	}
	return json.Unmarshal(b, p)
}
