package enums

type PetKind uint

const (
	Dog PetKind = iota + 1
	Cat
	Bird
	Fish
	Reptile
	Amphibian
	Rodent
	Rabbit
	Ferret
	Horse
	Turtle
	Snake
	Lizard
	Hedgehog
	MiniPig
	Insect
	Arachnid
	HermitCrab
)

func (petKind PetKind) String() string {
	switch petKind {
	case Dog:
		return "سگ"
	case Cat:
		return "گربه"
	case Bird:
		return "پرنده"
	case Fish:
		return "ماهی"
	case Reptile:
		return "خزنده"
	case Amphibian:
		return "دوزیست"
	case Rodent:
		return "جوندگان"
	case Rabbit:
		return "خرگوش"
	case Ferret:
		return "فرت / راسو"
	case Horse:
		return "اسب"
	case Turtle:
		return "لاک‌پشت"
	case Snake:
		return "مار"
	case Lizard:
		return "مارمولک"
	case Hedgehog:
		return "جوجه‌تیغی"
	case MiniPig:
		return "مینی‌پیگ / خوکچه خانگی"
	case Insect:
		return "حشره"
	case Arachnid:
		return "عنکبوت‌سانان"
	case HermitCrab:
		return "خرچنگ هرمت"
	}
	return ""
}

func GetAllPetKinds() []PetKind {
	return []PetKind{
		Dog,
		Cat,
		Bird,
		Fish,
		Reptile,
		Amphibian,
		Rodent,
		Rabbit,
		Ferret,
		Horse,
		Turtle,
		Snake,
		Lizard,
		Hedgehog,
		MiniPig,
		Insect,
		Arachnid,
		HermitCrab,
	}
}
