package enums

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
