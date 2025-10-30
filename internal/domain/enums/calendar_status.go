package enums

type CalendarStatus uint

const (
	Free CalendarStatus = iota + 1
	Booked
)

func (calenderStatus CalendarStatus) String() string {
	switch calenderStatus {
	case CalendarStatus(Free):
		return "free"
	case CalendarStatus(Booked):
		return "booked"
	}
	return ""
}

func GetAllCalenderStatus() []CalendarStatus {
	return []CalendarStatus{
		Free,
		Booked,
	}
}
