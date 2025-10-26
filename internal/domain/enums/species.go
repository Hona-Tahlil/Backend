package enums

type Species uint

const (
	German Species = iota + 1
)

func (species Species) String() string {
	switch species {
	case Species(German):
		return "ژرمن"
	}
	return ""
}

func GetAllSpecies() []Species {
	return []Species{
		German,
	}
}
