package enums

type Province uint

const (
	TehranProvince Province = iota + 1
	IsfahanProvince
)

func (province Province) String() string {
	switch province {
	case Province(TehranProvince):
		return "tehran"
	case Province(IsfahanProvince):
		return "isfahan"
	}
	return ""
}

func GetAllProvinces() []Province {
	return []Province{
		TehranProvince,
		IsfahanProvince,
	}
}
