package enums

type PetKind uint

const (
	Dog PetKind = iota + 1
	Cat
	Rabbit
)

func (petKind PetKind) String() string {
	switch petKind {
	case PetKind(Dog):
		return "dog"
	case PetKind(Cat):
		return "cat"
	case PetKind(Rabbit):
		return "rabbit"
	}
	return ""
}

func GetAllPetKinds() []PetKind {
	return []PetKind{
		Dog,
		Cat,
		Rabbit,
	}
}
