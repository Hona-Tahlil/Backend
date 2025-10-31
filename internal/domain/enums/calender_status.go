package enums

type CalenderStatus uint

const (
	Free CalenderStatus = iota + 1
	Booked
)

func (calenderStatus CalenderStatus) String() string {
	switch calenderStatus {
	case CalenderStatus(Free):
		return "free"
	case CalenderStatus(Booked):
		return "booked"
	}
	return ""
}

func GetAllCalenderStatus() []CalenderStatus {
	return []CalenderStatus{
		Free,
		Booked,
	}
}
