package enums

import "fmt"

type Slot uint

const TotalSlots = 48

func (slot Slot) String() string {
	if slot < 1 || slot > TotalSlots {
		return ""
	}
	idx := int(slot - 1)
	startMin := idx * 30
	endMin := startMin + 30

	return fmt.Sprintf("%s-%s", fmtTime(startMin), fmtTime(endMin))
}

func fmtTime(min int) string {
	if min == 24*60 {
		return "24:00"
	}
	h := min / 60
	m := min % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

func GetAllSlots() []Slot {
	slots := make([]Slot, TotalSlots)
	for i := 1; i <= TotalSlots; i++ {
		slots[i-1] = Slot(i)
	}
	return slots
}
