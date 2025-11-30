package enums

type Species uint

const (
	// Dogs
	GermanShepherd Species = iota + 1
	GoldenRetriever
	Bulldog
	Poodle
	Labrador
	Husky
	Beagle
	Rottweiler
	Chihuahua
	ShihTzu
	Boxer
	Pomeranian
	Dachshund
	Corgi

	// Cats
	Persian
	Siamese
	MaineCoon
	BritishShorthair
	Sphynx
	Bengal
	Ragdoll
	ScottishFold

	// Birds
	Parrot
	Canary
	Budgie
	Cockatiel
	Finch
	Lovebird
	Macaw
	Cockatoo

	// Fish
	Goldfish
	Betta
	Guppy
	NeonTetra
	Discus
	Angelfish

	// Rodents
	SyrianHamster
	DwarfHamster
	GuineaPig
	Mouse
	Rat
	Gerbil
	Chinchilla

	// Rabbits
	HollandLop
	MiniRex
	Lionhead
	NetherlandDwarf

	// Reptiles – Turtles / Snakes / Lizards
	RedEaredSlider
	SulcataTortoise
	CornSnake
	BallPython
	LeopardGecko
	BeardedDragon
	CrestedGecko
	Iguana

	// Amphibians
	Axolotl
	TreeFrog
	PacmanFrog

	// Ferret
	StandardFerret

	// Horses
	MiniHorse

	// Hedgehog
	AfricanPygmyHedgehog

	// Mini Pig
	PotBelliedPig

	// Insects
	StickInsect
	PrayingMantis

	// Arachnids
	Tarantula

	// Hermit Crab
	CaribbeanHermitCrab

	//Other
	Other
)

func (species Species) String() string {
	switch species {

	// Dogs
	case GermanShepherd:
		return "ژرمن شپرد"
	case GoldenRetriever:
		return "گلدن رتریور"
	case Bulldog:
		return "بولداگ"
	case Poodle:
		return "پودل"
	case Labrador:
		return "لابرادور"
	case Husky:
		return "هاسکی"
	case Beagle:
		return "بیگل"
	case Rottweiler:
		return "روتوایلر"
	case Chihuahua:
		return "چیهواوا"
	case ShihTzu:
		return "شیه‌تزو"
	case Boxer:
		return "باکسر"
	case Pomeranian:
		return "پامرانیان"
	case Dachshund:
		return "داشهوند"
	case Corgi:
		return "کورگی"

	// Cats
	case Persian:
		return "پرشین"
	case Siamese:
		return "سیامی"
	case MaineCoon:
		return "مین کون"
	case BritishShorthair:
		return "بریتیش شورتهیر"
	case Sphynx:
		return "اسفنکس"
	case Bengal:
		return "بنگال"
	case Ragdoll:
		return "رagdoll"
	case ScottishFold:
		return "اسکاتیش فولد"

	// Birds
	case Parrot:
		return "طوطی"
	case Canary:
		return "قناری"
	case Budgie:
		return "مرغ عشق"
	case Cockatiel:
		return "کاکاتیل"
	case Finch:
		return "فنچ"
	case Lovebird:
		return "عاشق‌پرنده"
	case Macaw:
		return "ماکائو"
	case Cockatoo:
		return "کاکاتو"

	// Fish
	case Goldfish:
		return "گلدفیش"
	case Betta:
		return "بِتا"
	case Guppy:
		return "گوپی"
	case NeonTetra:
		return "نئون تترا"
	case Discus:
		return "دیسکاس"
	case Angelfish:
		return "آنجل فیش"

	// Rodents
	case SyrianHamster:
		return "همستر سوری"
	case DwarfHamster:
		return "همستر کوتوله"
	case GuineaPig:
		return "خوکچه هندی"
	case Mouse:
		return "موش"
	case Rat:
		return "رت"
	case Gerbil:
		return "جربیل"
	case Chinchilla:
		return "چینچیلا"

	// Rabbits
	case HollandLop:
		return "هلند لاپ"
	case MiniRex:
		return "مینی رِکس"
	case Lionhead:
		return "لاین‌هد"
	case NetherlandDwarf:
		return "نِترلند دورف"

	// Reptiles
	case RedEaredSlider:
		return "لاک‌پشت گوش‌قرمز"
	case SulcataTortoise:
		return "لاک‌پشت سولکاتا"
	case CornSnake:
		return "مار ذرت"
	case BallPython:
		return "پایتون توپی"
	case LeopardGecko:
		return "گکو پلنگی"
	case BeardedDragon:
		return "اژدهای ریش‌دار"
	case CrestedGecko:
		return "گکوی کرستد"
	case Iguana:
		return "ایگوانا"

	// Amphibians
	case Axolotl:
		return "آکسولوتل"
	case TreeFrog:
		return "قورباغه درختی"
	case PacmanFrog:
		return "قورباغه پاکمن"

	// Ferret
	case StandardFerret:
		return "فرت"

	// Horses
	case MiniHorse:
		return "اسب مینیاتوری"

	// Hedgehog
	case AfricanPygmyHedgehog:
		return "جوجۀ تیغی آفریقایی"

	// Mini Pig
	case PotBelliedPig:
		return "مینی‌پیگ"

	// Insects
	case StickInsect:
		return "حشره چوبی"
	case PrayingMantis:
		return "مانتیس"

	// Arachnids
	case Tarantula:
		return "تارانتولا"

	// Hermit Crab
	case CaribbeanHermitCrab:
		return "خرچنگ هرمت کارائیبی"

	default:
		return "دیگر نژاد ها / نمی دانم"
	}
}

func GetAllSpecies() []Species {
	return []Species{
		// Dogs
		GermanShepherd, GoldenRetriever, Bulldog, Poodle, Labrador, Husky,
		Beagle, Rottweiler, Chihuahua, ShihTzu, Boxer, Pomeranian, Dachshund, Corgi,

		// Cats
		Persian, Siamese, MaineCoon, BritishShorthair, Sphynx, Bengal, Ragdoll, ScottishFold,

		// Birds
		Parrot, Canary, Budgie, Cockatiel, Finch, Lovebird, Macaw, Cockatoo,

		// Fish
		Goldfish, Betta, Guppy, NeonTetra, Discus, Angelfish,

		// Rodents
		SyrianHamster, DwarfHamster, GuineaPig, Mouse, Rat, Gerbil, Chinchilla,

		// Rabbits
		HollandLop, MiniRex, Lionhead, NetherlandDwarf,

		// Reptiles
		RedEaredSlider, SulcataTortoise, CornSnake, BallPython, LeopardGecko,
		BeardedDragon, CrestedGecko, Iguana,

		// Amphibians
		Axolotl, TreeFrog, PacmanFrog,

		// Ferret
		StandardFerret,

		// Horse
		MiniHorse,

		// Hedgehog
		AfricanPygmyHedgehog,

		// Mini Pig
		PotBelliedPig,

		// Insects
		StickInsect, PrayingMantis,

		// Arachnids
		Tarantula,

		// Hermit Crab
		CaribbeanHermitCrab,

		//Other
		Other,
	}
}

func GetSpeciesByKind(kind PetKind) []Species {
	switch kind {

	case Dog:
		return []Species{
			GermanShepherd,
			GoldenRetriever,
			Bulldog,
			Poodle,
			Labrador,
			Husky,
			Beagle,
			Rottweiler,
			Chihuahua,
			ShihTzu,
			Boxer,
			Pomeranian,
			Dachshund,
			Corgi,
		}

	case Cat:
		return []Species{
			Persian,
			Siamese,
			MaineCoon,
			BritishShorthair,
			Sphynx,
			Bengal,
			Ragdoll,
			ScottishFold,
		}

	case Bird:
		return []Species{
			Parrot,
			Canary,
			Budgie,
			Cockatiel,
			Finch,
			Lovebird,
			Macaw,
			Cockatoo,
		}

	case Fish:
		return []Species{
			Goldfish,
			Betta,
			Guppy,
			NeonTetra,
			Discus,
			Angelfish,
		}

	case Rodent:
		return []Species{
			SyrianHamster,
			DwarfHamster,
			GuineaPig,
			Mouse,
			Rat,
			Gerbil,
			Chinchilla,
		}

	case Rabbit:
		return []Species{
			HollandLop,
			MiniRex,
			Lionhead,
			NetherlandDwarf,
		}

	case Reptile:
		return []Species{
			RedEaredSlider,
			SulcataTortoise,
			CornSnake,
			BallPython,
			LeopardGecko,
			BeardedDragon,
			CrestedGecko,
			Iguana,
		}

	case Amphibian:
		return []Species{
			Axolotl,
			TreeFrog,
			PacmanFrog,
		}

	case Ferret:
		return []Species{
			StandardFerret,
		}

	case Horse:
		return []Species{
			MiniHorse,
		}

	case Hedgehog:
		return []Species{
			AfricanPygmyHedgehog,
		}

	case MiniPig:
		return []Species{
			PotBelliedPig,
		}

	case Insect:
		return []Species{
			StickInsect,
			PrayingMantis,
		}

	case Arachnid:
		return []Species{
			Tarantula,
		}

	case HermitCrab:
		return []Species{
			CaribbeanHermitCrab,
		}

	default:
		return []Species{}
	}
}
