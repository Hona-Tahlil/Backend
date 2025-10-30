package enums

type Gender uint

const (
	Male Gender = iota + 1
	Female
	// Unknown
)

func (gender Gender) String() string {
	switch gender {
	case Gender(Male):
		return "male"
	case Gender(Female):
		return "female"
		// case Gender(Unknown):
		// 	return "unknown"
	}
	return ""
}

func GetAllGenders() []Gender {
	return []Gender{
		Male,
		Female,
		// Unknown,
	}
}
