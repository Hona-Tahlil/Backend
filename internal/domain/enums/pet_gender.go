package enums

type PetGender uint

const (
	MalePet PetGender = iota + 1
	FemalePet
	Unknown
)

func (petGender PetGender) String() string {
	switch petGender {
	case PetGender(MalePet):
		return "نر"
	case PetGender(FemalePet):
		return "ماده"
	case PetGender(Unknown):
		return "نامشخص"
	}
	return ""
}

func GetAllPetGenders() []PetGender {
	return []PetGender{
		MalePet,
		FemalePet,
		Unknown,
	}
}
