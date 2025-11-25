package enums

type City uint

const (
	TehranCity City = iota + 1
	IsfahanCity
)

func (city City) String() string {
	switch city {
	case City(TehranCity):
		return "tehran"
	case City(IsfahanCity):
		return "isfahan"
	}
	return ""
}

func GetAllCities() []City {
	return []City{
		TehranCity,
		IsfahanCity,
	}
}
