package seeder

import (
	"fmt"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/ports"

	"gorm.io/gorm"
)

var provinceWithCities = map[enums.Province][]enums.City{
	enums.EastAzerbaijan: {
		enums.TabrizCity,
		enums.MaraghehCity,
		enums.MiyanehCity,
		enums.ShabestarCity,
		enums.MarandCity,
		enums.JolfaCity,
		enums.SarabCity,
		enums.HadishahrCity,
		enums.BonabCity,
		enums.TasujCity,
		enums.AharCity,
		enums.HashtrudCity,
		enums.MalekanCity,
		enums.Bostan_AbadCity,
		enums.VarzaqanCity,
		enums.OskuCity,
		enums.MamqanCity,
		enums.SofianCity,
		enums.IlkhchiCity,
		enums.KhosrowshahrCity,
		enums.BasmenjCity,
		enums.SahandCity,
	},
	enums.WestAzerbaijan: {
		enums.UrmiaCity,
		enums.NaqadehCity,
		enums.MakuCity,
		enums.TakabCity,
		enums.KhoyCity,
		enums.MahabadCity,
		enums.SardashtCity,
		enums.ChaldoranCity,
		enums.BukanCity,
		enums.MiandoabCity,
		enums.SalmasCity,
		enums.Shahin_DezhCity,
		enums.PiranshahrCity,
		enums.OshnaviehCity,
		enums.PoldashtCity,
	},
	enums.Ardabil: {
		enums.ArdabilCity,
		enums.ParsabadCity,
		enums.KhalkhalCity,
		enums.Meshgin_ShahrCity,
		enums.NaminCity,
		enums.NirCity,
		enums.GarmiCity,
	},
	enums.Isfahan: {
		enums.IsfahanCity,
		enums.FalavarjanCity,
		enums.GolpayeganCity,
		enums.DehaghanCity,
		enums.NatanzCity,
		enums.TiranCity,
		enums.KashanCity,
		enums.ArdestanCity,
		enums.SemiromCity,
		enums.DorchehCity,
		enums.KouhpayehCity,
		enums.MobarekehCity,
		enums.Shahr_e_RezaCity,
		enums.Khomeyni_ShahrCity,
		enums.NajafabadCity,
		enums.Zarrin_ShahrCity,
		enums.Aran_o_BidgolCity,
		enums.Baq_e_BahadoranCity,
		enums.KhvansarCity,
		enums.AlavijehCity,
		enums.AsgaranCity,
		enums.HajiabadCity,
		enums.TudeshkCity,
		enums.VarzanehCity,
	},
	enums.Ilam: {
		enums.IlamCity,
		enums.MehranCity,
		enums.DehloranCity,
		enums.AbdananCity,
		enums.Darreh_ShahrCity,
		enums.EyvanCity,
		enums.SarablehCity,
	},
	enums.Bushehr: {
		enums.BushehrCity,
		enums.DayyerCity,
		enums.KanganCity,
		enums.GenavehCity,
		enums.KhourmojCity,
		enums.AhramCity,
		enums.BorazjanCity,
		enums.JamCity,
		enums.KakiCity,
		enums.AsaluyehCity,
	},
	enums.Tehran: {
		enums.TehranCity,
		enums.VaraminCity,
		enums.FiroozkoohCity,
		enums.ReyCity,
		enums.DamavandCity,
		enums.EslamshahrCity,
		enums.RodehenCity,
		enums.LavasanCity,
		enums.BoomhenCity,
		enums.TajrishCity,
		enums.FashamCity,
		enums.KahrizakCity,
		enums.PakdashtCity,
		enums.ChahardangehCity,
		enums.Sharif_AbadCity,
		enums.QarchakCity,
		enums.BaqershahrCity,
		enums.ShahriarCity,
		enums.Robat_KarimCity,
		enums.QodsCity,
		enums.MalardCity,
	},
	enums.ChaharmahalBakhtiari: {
		enums.ShahrekordCity,
		enums.FarsanCity,
		enums.BorujenCity,
		enums.ChelgardCity,
		enums.ArdalCity,
		enums.LordeganCity,
	},
	enums.SouthKhorasan: {
		enums.QaenCity,
		enums.FerdowsCity,
		enums.BirjandCity,
		enums.NehbandanCity,
		enums.SarbishchCity,
		enums.TabasCity,
	},
	enums.RazaviKhorasan: {
		enums.MashhadCity,
		enums.NeyshaburCity,
		enums.SabzevarCity,
		enums.KashmarCity,
		enums.GonabadCity,
		enums.Torbat_HeydariehCity,
		enums.KhafCity,
		enums.Torbat_e_JamCity,
		enums.TaybadCity,
		enums.QuchanCity,
		enums.SarakhsCity,
		enums.FarimanCity,
		enums.ChenaranCity,
		enums.DargazCity,
		enums.ToroqbehCity,
	},
	enums.NorthKhorasan: {
		enums.BojnurdCity,
		enums.EsfarayenCity,
		enums.JajarmCity,
		enums.ShirvanCity,
		enums.AshkhanehCity,
	},
	enums.Khuzestan: {
		enums.AhvazCity,
		enums.ShushCity,
		enums.AbadanCity,
		enums.KhorramshahrCity,
		enums.Masjed_SoleymanCity,
		enums.IzehCity,
		enums.ShushtarCity,
		enums.AndimeshkCity,
		enums.SusangerdCity,
		enums.HoveyzehCity,
		enums.DezfulCity,
		enums.ShadeganCity,
		enums.Bandar_MahshahrCity,
		enums.Bandar_Imam_KhomeiniCity,
		enums.BehbahanCity,
		enums.RamhormozCity,
		enums.Bagh_e_MalekCity,
		enums.HendijanCity,
		enums.LaliCity,
		enums.RamshirCity,
		enums.HamidiyehCity,
		enums.MollasaniCity,
	},
	enums.Zanjan: {
		enums.ZanjanCity,
		enums.AbharCity,
		enums.KhodabandehCity,
		enums.MahneshanCity,
		enums.KhorramdarrehCity,
		enums.Ab_BarCity,
		enums.QeydarCity,
	},
	enums.Semnan: {
		enums.SemnanCity,
		enums.ShahroodCity,
		enums.GarmsarCity,
		enums.EyvankiCity,
		enums.DamghanCity,
		enums.BastamCity,
	},
	enums.SistanBaluchestan: {
		enums.ZahedanCity,
		enums.ChabaharCity,
		enums.KhashCity,
		enums.SaravanCity,
		enums.ZabolCity,
		enums.SarbazCity,
		enums.MirjavehCity,
	},
	enums.Fars: {
		enums.ShirazCity,
		enums.EghlidCity,
		enums.DarabCity,
		enums.FasaCity,
		enums.MarvdashtCity,
		enums.AbadehCity,
		enums.KazerunCity,
		enums.SepidanCity,
		enums.LarCity,
		enums.FiroozabadCity,
		enums.JahromCity,
		enums.EstahbanCity,
		enums.LamerdCity,
		enums.MohrCity,
		enums.ArdakanCity,
		enums.SafashahrCity,
		enums.ArsanjanCity,
		enums.SurianCity,
		enums.FarashbandCity,
		enums.SarvestanCity,
		enums.ZarghanCity,
		enums.KavarCity,
		enums.BavanatCity,
		enums.KharamehCity,
		enums.KhonjCity,
	},
	enums.Qazvin: {
		enums.QazvinCity,
		enums.TakestanCity,
		enums.AbyekCity,
		enums.Boin_ZahraCity,
	},
	enums.Qom: {
		enums.QomCity,
		enums.QanavatCity,
		enums.JafariyehCity,
		enums.KahakCity,
		enums.DastjerdCity,
		enums.SalfcheganCity,
	},
	enums.Alborz: {
		enums.KarajCity,
		enums.TaleqanCity,
		enums.NazarabadCity,
		enums.EshtehardCity,
		enums.HashtgerdCity,
		enums.MahdashtCity,
	},
	enums.Kurdistan: {
		enums.SanandajCity,
		enums.BanehCity,
		enums.BijarCity,
		enums.SaqqezCity,
		enums.QorvehCity,
		enums.MarivanCity,
		enums.Solvat_AbadCity,
		enums.Hasan_AbadCity,
	},
	enums.Kerman: {
		enums.KermanCity,
		enums.RavarCity,
		enums.AnarCity,
		enums.KuhbananCity,
		enums.RafsanjanCity,
		enums.BaftCity,
		enums.SirjanCity,
		enums.KahnujCity,
		enums.ZarandCity,
		enums.BamCity,
		enums.JiroftCity,
		enums.BardsirCity,
	},
	enums.Kermanshah: {
		enums.KermanshahCity,
		enums.Islamabad_GharbCity,
		enums.KangavarCity,
		enums.SonqorCity,
		enums.Qasr_e_ShirinCity,
		enums.HarsinCity,
		enums.PavehCity,
		enums.JavanrudCity,
	},
	enums.KohgiluyehBoyerAhmad: {
		enums.YasujCity,
		enums.GachsaranCity,
		enums.DogonbadanCity,
		enums.SisakhtCity,
		enums.DehdashtCity,
	},
	enums.Golestan: {
		enums.GorganCity,
		enums.Aq_QalaCity,
		enums.Gonbad_e_KavusCity,
		enums.Ali_Abad_e_KatulCity,
		enums.KordkuyCity,
		enums.KalalehCity,
		enums.AzadshahrCity,
		enums.RamiyanCity,
	},
	enums.Gilan: {
		enums.RashtCity,
		enums.ManjilCity,
		enums.LangarudCity,
		enums.TaleshCity,
		enums.AstaraCity,
		enums.MasoulehCity,
		enums.RoudbarCity,
		enums.FumanCity,
		enums.Sowme_eh_SaraCity,
		enums.HashtparCity,
		enums.MasalCity,
		enums.ShaftCity,
		enums.AmlashCity,
		enums.LahijanCity,
	},
	enums.Lorestan: {
		enums.KhorramabadCity,
		enums.DorudCity,
		enums.AligudarzCity,
		enums.AznaCity,
		enums.Noor_AbadCity,
		enums.KuhdashtCity,
		enums.AleshtarCity,
	},
	enums.Mazandaran: {
		enums.SariCity,
		enums.AmolCity,
		enums.BabolCity,
		enums.BabolsarCity,
		enums.BehshahrCity,
		enums.TonekabonCity,
		enums.JuybarCity,
		enums.ChalusCity,
		enums.RamsarCity,
		enums.Qaem_ShahrCity,
		enums.NekaCity,
		enums.NoorCity,
		enums.BaladehCity,
		enums.NowshahrCity,
		enums.MahmudabadCity,
	},
	enums.Markazi: {
		enums.ArakCity,
		enums.AshtianCity,
		enums.TafreshCity,
		enums.KhomeinCity,
		enums.DelijanCity,
		enums.SavehCity,
		enums.MahallatCity,
		enums.ShazandCity,
	},
	enums.Hormozgan: {
		enums.Bandar_AbbasCity,
		enums.QeshmCity,
		enums.KishCity,
		enums.Bandar_LengehCity,
		enums.BastakCity,
		enums.DehbarzCity,
		enums.MinabCity,
		enums.Bandar_JaskCity,
		enums.Bandar_KhamirCity,
	},
	enums.Hamadan: {
		enums.HamedanCity,
		enums.MalayerCity,
		enums.NahavandCity,
		enums.RazanCity,
		enums.AsadabadCity,
		enums.BaharCity,
	},
	enums.Yazd: {
		enums.YazdCity,
		enums.TaftCity,
		enums.Ardakan_YazdCity,
		enums.AbarkouhCity,
		enums.MeybodCity,
		enums.BafqCity,
		enums.MehrizCity,
		enums.AshkezarCity,
		enums.HarratCity,
		enums.KhezrabadCity,
		enums.ZarchCity,
	},
}

type AddressSeeder struct {
	unitOfwork ports.UnitOfWork
	db         *gorm.DB
}

func NewAddressSeeder(
	unitOfwork ports.UnitOfWork,
	db *gorm.DB,
) *AddressSeeder {
	return &AddressSeeder{
		unitOfwork: unitOfwork,
		db:         db,
	}
}

func (seeder *AddressSeeder) SeedProvincesAndCities() {
	for provinceName, cityNames := range provinceWithCities {
		province, err := seeder.unitOfwork.Factory().ProvinceRepository().FindProvinceByName(provinceName)
		if err != nil {
			panic(err)
		}
		if province == nil {
			province = &entities.Province{
				Name: provinceName,
			}
			err := seeder.unitOfwork.Factory().ProvinceRepository().CreateProvince(province)
			if err != nil {
				panic(fmt.Errorf("error creating province %s: %w", provinceName, err))
			}

		}

		for _, cityName := range cityNames {
			var city *entities.City
			for _, pCity := range province.Cities {
				if pCity.Name == cityName {
					city = &pCity
					break
				}
			}
			if city != nil {
				continue
			}
			city = &entities.City{
				Name:       cityName,
				ProvinceID: province.ID,
			}
			err = seeder.unitOfwork.Factory().CityRepository().CreateCity(city)
			if err != nil {
				panic(fmt.Errorf("error creating city %s: %w", cityName, err))
			}
		}
	}
}
