package enums

type Province uint

const (
	Alborz Province = iota + 1 // Starts at 1
	Ardabil
	Bushehr
	ChaharmahalBakhtiari
	EastAzerbaijan
	Fars
	Gilan
	Golestan
	Hamadan
	Hormozgan
	Ilam
	Isfahan
	Kerman
	Kermanshah
	Khuzestan
	KohgiluyehBoyerAhmad
	Kurdistan
	Lorestan
	Markazi
	Mazandaran
	NorthKhorasan
	Qazvin
	Qom
	RazaviKhorasan
	Semnan
	SistanBaluchestan
	SouthKhorasan
	Tehran
	WestAzerbaijan
	Yazd
	Zanjan
)

func (p Province) String() string {
	if int(p) >= len(FarsiNames) || p < 1 {
		return "ProvinceUnknown"
	}
	return FarsiNames[p]
}

func GetAllProvinces() []Province {
	provinces := make([]Province, 0, 31)
	for i := Alborz; i < 32; i++ {
		provinces[i-1] = i
	}
	return provinces
}

var FarsiNames = map[Province]string{
	EastAzerbaijan:       "آذربايجان شرقی",
	WestAzerbaijan:       "آذربايجان غربی",
	Ardabil:              "اردبيل",
	Isfahan:              "اصفهان",
	Ilam:                 "ايلام",
	Bushehr:              "بوشهر",
	Tehran:               "تهران",
	ChaharmahalBakhtiari: "چهارمحال بختیاری",
	SouthKhorasan:        "خراسان جنوبی",
	RazaviKhorasan:       "خراسان رضوی",
	NorthKhorasan:        "خراسان شمالی",
	Khuzestan:            "خوزستان",
	Zanjan:               "زنجان",
	Semnan:               "سمنان",
	SistanBaluchestan:    "سيستان و بلوچستان",
	Fars:                 "فارس",
	Qazvin:               "قزوين",
	Qom:                  "قم",
	Alborz:               "البرز",
	Kurdistan:            "كردستان",
	Kerman:               "کرمان",
	Kermanshah:           "كرمانشاه",
	KohgiluyehBoyerAhmad: "کهگیلویه و بويراحمد",
	Golestan:             "گلستان",
	Gilan:                "گيلان",
	Lorestan:             "لرستان",
	Mazandaran:           "مازندران",
	Markazi:              "مرکزی",
	Hormozgan:            "هرمزگان",
	Hamadan:              "همدان",
	Yazd:                 "يزد",
}
