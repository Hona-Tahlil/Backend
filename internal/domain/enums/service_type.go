package enums

type ServiceType uint

const (
	Walking ServiceType = iota + 1
	Watching
)

func (serviceType ServiceType) String() string {
	switch serviceType {
	case ServiceType(Walking):
		return "walking"
	case ServiceType(Watching):
		return "watching"
	}
	return ""
}

func GetAllServiceTypes() []ServiceType {
	return []ServiceType{
		Walking,
		Watching,
	}
}
