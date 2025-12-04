package enums

type PetKind uint

const (
	Dog PetKind = iota + 1
	Cat
	Bird
	Rodent
)

func (petKind PetKind) String() string {
	switch petKind {
	case Dog:
		return "سگ"
	case Cat:
		return "گربه"
	case Bird:
		return "پرنده"
	case Rodent:
		return "جونده"
	}
	return ""
}

func GetAllPetKinds() []PetKind {
	return []PetKind{
		Dog,
		Cat,
		Bird,
		Rodent,
	}
}
