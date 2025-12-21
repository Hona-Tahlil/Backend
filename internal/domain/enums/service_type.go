package enums

type ServiceType uint

const (
	Walking ServiceType = iota + 1
	Watching
)

func (serviceType ServiceType) String() string {
	switch serviceType {
	case ServiceType(Walking):
		return "پیاده روی"
	case ServiceType(Watching):
		return "نگهداری"
	}
	return ""
}

func GetAllServiceTypes() []ServiceType {
	return []ServiceType{
		Walking,
		Watching,
	}
}
