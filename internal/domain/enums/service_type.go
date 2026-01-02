package enums

type ServiceType uint

const (
	Walking ServiceType = iota + 1
	Training
	Watching
	MedicalCare
)

func (serviceType ServiceType) String() string {
	switch serviceType {
	case ServiceType(Walking):
		return "پیاده روی"
	case ServiceType(Training):
		return "آموزش"
	case ServiceType(Watching):
		return "نگهداری"
	case ServiceType(MedicalCare):
		return "مراقبت های پزشکی"
	}
	return ""
}

func GetAllServiceTypes() []ServiceType {
	return []ServiceType{
		Walking,
		Training,
		Watching,
		MedicalCare,
	}
}
