package seeder

import (
	"fmt"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"

	"gorm.io/gorm"
)

type DatabaseSeeder struct {
	db *gorm.DB
}

func NewDatabaseSeeder(db *gorm.DB) *DatabaseSeeder {
	return &DatabaseSeeder{db: db}
}

// SeedAll runs all seeders in correct order with default counts
func (s *DatabaseSeeder) SeedAll() error {
	log.Println("🌱 Starting database seeding...")

	// if err := s.SeedUsers(50); err != nil {
	// 	return err
	// }
	if err := s.SeedUsers(50); err != nil {
		log.Println("SeedUsers error:", err)
		return err
	}

	if err := s.SeedPets(150); err != nil {
		return err
	}

	// مهم: Address باید قبل از سرچ آماده باشد
	if err := s.SeedAddressesForUsers(); err != nil {
		return err
	}

	if err := s.SeedPetSitters(20); err != nil {
		return err
	}

	// سرویس‌ها باید بر اساس pet_sitter های واقعی ساخته شوند
	if err := s.SeedServicesForPetSitters(); err != nil {
		return err
	}

	// برای فیلتر date/slot
	if err := s.SeedCalendarSlotsForPetSitters(250); err != nil {
		return err
	}

	// برای rate و commentsCount
	if err := s.SeedCommentsForPetSitters(400); err != nil {
		return err
	}

	log.Println("✅ Database seeding completed!")
	return nil
}

// SeedUsers seeds user data
func (s *DatabaseSeeder) SeedUsers(count int) error {
	seeder := NewUserSeeder(s.db)
	return seeder.Seed(count)
}

// SeedPets seeds pet data
func (s *DatabaseSeeder) SeedPets(count int) error {
	seeder := NewPetSeeder(s.db)
	return seeder.Seed(count)
}

// SeedPetSitters seeds pet sitter data
func (s *DatabaseSeeder) SeedPetSitters(count int) error {
	seeder := NewPetSitterSeeder(s.db)
	return seeder.Seed(count)
}

// SeedServices seeds service data
func (s *DatabaseSeeder) SeedServices(count int) error {
	return ServiceSeeder(s.db)
}

// SeedRequests seeds request data
func (s *DatabaseSeeder) SeedRequests(count int) error {
	// seeder := NewRequestSeeder(s.db)
	// return seeder.Seed(count)
	return nil
}

// SeedWallets seeds wallet data
func (s *DatabaseSeeder) SeedWallets(count int) error {
	log.Printf("💰 Seeding %d wallets...", count)
	log.Println("⚠️  Wallet seeding not implemented yet")
	return nil
}

// SeedAddresses seeds address data
func (s *DatabaseSeeder) SeedAddresses(count int) error {
	log.Printf("📍 Seeding %d addresses...", count)
	log.Println("⚠️  Address seeding not implemented yet")
	return nil
}

// SeedChats seeds chat data
func (s *DatabaseSeeder) SeedChats(count int) error {
	log.Printf("💬 Seeding %d chats...", count)
	log.Println("⚠️  Chat seeding not implemented yet")
	return nil
}

func (s *DatabaseSeeder) ClearAll() {
	log.Println("🗑️  Clearing all tables...")

	s.db.Migrator().DropTable(
		"user_roles",
		&entities.User{},
		&entities.Role{},
		&entities.Permission{},
		&entities.Wallet{},
		&entities.Request{},
		&entities.CalendarSlot{},
		&entities.Pet{},
		&entities.PetSitter{},
		&entities.Service{},
		&entities.ChatRoom{},
		&entities.Comment{},
		&entities.Address{},
		&entities.ChatMessage{},
		&entities.Transaction{},
		&entities.Transfer{},
	)

}

func (s *DatabaseSeeder) SeedServicesForPetSitters() error {
	// اگر service وجود دارد، skip
	if s.db.First(&entities.Service{}).Error == nil {
		log.Println("✓ Services already seeded, skipping...")
		return nil
	}

	var sitters []entities.PetSitter
	if err := s.db.Find(&sitters).Error; err != nil {
		return err
	}
	if len(sitters) == 0 {
		return fmt.Errorf("no pet sitters found, seed pet sitters first")
	}

	allTypes := enums.GetAllServiceTypes() // اگر نداری، لیست دستی بده
	services := make([]entities.Service, 0, len(sitters)*2)

	for _, ps := range sitters {
		// چند type تصادفی ولی بدون تکرار
		types := pickUniqueServiceTypes(allTypes, rand.Intn(min(3, len(allTypes)))+1)

		for _, t := range types {
			desc := faker.Sentence()
			svc := entities.Service{
				PetSitterID: ps.ID,
				Type:        t,
				Price:       uint(rand.Intn(90)+10) * 1000, // 10k تا 100k
				Description: &desc,
				Kind:        "petSitter",
			}
			services = append(services, svc)
		}
	}

	if err := s.db.CreateInBatches(services, 200).Error; err != nil {
		return fmt.Errorf("seed services: %w", err)
	}

	log.Printf("✓ Successfully seeded %d services", len(services))
	return nil
}

func (s *DatabaseSeeder) SeedAddressesForUsers() error {
	if s.db.First(&entities.Address{}).Error == nil {
		log.Println("✓ Addresses already seeded, skipping...")
		return nil
	}

	var users []entities.User
	if err := s.db.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("no users found, seed users first")
	}

	addresses := make([]entities.Address, 0, len(users))

	for _, u := range users {
		addr := entities.Address{
			Province: randomProvince(),
			City:     randomCity(),
			// StreetAddress: faker.StreetAddress(),
			HouseNumber: uint(rand.Intn(200) + 1),
			Unit:        uint(rand.Intn(10) + 1),
			PostalCode:  nil,
			Refer:       u.ID,
			Type:        "User",
		}
		addresses = append(addresses, addr)
	}

	if err := s.db.CreateInBatches(addresses, 200).Error; err != nil {
		return fmt.Errorf("seed addresses: %w", err)
	}

	log.Printf("✓ Successfully seeded %d addresses", len(addresses))
	return nil
}

func (s *DatabaseSeeder) SeedCalendarSlotsForPetSitters(count int) error {
	if s.db.First(&entities.CalendarSlot{}).Error == nil {
		log.Println("✓ Calendar slots already seeded, skipping...")
		return nil
	}

	var sitters []entities.PetSitter
	if err := s.db.Find(&sitters).Error; err != nil {
		return err
	}
	if len(sitters) == 0 {
		return fmt.Errorf("no pet sitters found")
	}

	slots := make([]entities.CalendarSlot, 0, count)

	for i := 0; i < count; i++ {
		ps := sitters[rand.Intn(len(sitters))]

		start := rand.Intn(10)
		end := start + rand.Intn(4)
		slotArr := make([]enums.Slot, 0, end-start+1)
		for k := start; k <= end; k++ {
			slotArr = append(slotArr, enums.Slot(k))
		}

		cs := entities.CalendarSlot{
			PetSitterID: &ps.ID,
			Date:        time.Now().AddDate(0, 0, rand.Intn(14)),
			Slots:       slotArr,
			Status:      enums.Free, // اگر enum متفاوت است اصلاح کن
		}
		slots = append(slots, cs)
	}

	if err := s.db.CreateInBatches(slots, 200).Error; err != nil {
		return fmt.Errorf("seed calendar slots: %w", err)
	}

	log.Printf("✓ Successfully seeded %d calendar slots", len(slots))
	return nil
}

func (s *DatabaseSeeder) SeedCommentsForPetSitters(count int) error {
	if s.db.First(&entities.Comment{}).Error == nil {
		log.Println("✓ Comments already seeded, skipping...")
		return nil
	}

	// باید request های واقعی داشته باشیم
	var requests []entities.Request
	if err := s.db.Select("id").Find(&requests).Error; err != nil {
		return fmt.Errorf("fetch requests: %w", err)
	}
	if len(requests) == 0 {
		log.Println("⚠️  No requests found, skipping comment seeding (comments require valid request_id)")
		return nil
	}

	var sitters []entities.PetSitter
	if err := s.db.Find(&sitters).Error; err != nil {
		return err
	}
	if len(sitters) == 0 {
		return fmt.Errorf("no pet sitters found")
	}

	var users []entities.User
	if err := s.db.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("no users found")
	}

	comments := make([]entities.Comment, 0, count)

	for i := 0; i < count; i++ {
		ps := sitters[rand.Intn(len(sitters))]
		u := users[rand.Intn(len(users))]
		rq := requests[rand.Intn(len(requests))] // ✅ واقعی

		txt := faker.Sentence()
		c := entities.Comment{
			UserID:      u.ID,
			PetSitterID: ps.ID,
			RequestID:   rq.ID, // ✅ FK درست
			Text:        &txt,
			Rating:      uint(rand.Intn(5) + 1),
		}
		comments = append(comments, c)
	}

	if err := s.db.CreateInBatches(comments, 300).Error; err != nil {
		return fmt.Errorf("seed comments: %w", err)
	}

	log.Printf("✓ Successfully seeded %d comments", len(comments))
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func pickUniqueServiceTypes(all []enums.ServiceType, n int) []enums.ServiceType {
	if n <= 0 {
		return nil
	}
	if n >= len(all) {
		out := make([]enums.ServiceType, len(all))
		copy(out, all)
		return out
	}

	perm := rand.Perm(len(all))
	out := make([]enums.ServiceType, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, all[perm[i]])
	}
	return out
}

func randomPetKinds() []enums.PetKind {
	all := enums.GetAllPetKinds() // اگر داری
	if len(all) == 0 {
		return []enums.PetKind{}
	}

	n := rand.Intn(min(3, len(all))) + 1
	perm := rand.Perm(len(all))
	out := make([]enums.PetKind, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, all[perm[i]])
	}
	return out
}

// اگر GetAllProvince/GetAllCity نداری، مثل PetKinds عمل کن یا دستی لیست بده
// func randomProvince() enums.Province {
// 	all := enums.GetAllProvinces()
// 	return all[rand.Intn(len(all))]
// }

// func randomCity() enums.City {
// 	all := enums.GetAllCities()
// 	return all[rand.Intn(len(all))]
// }

func randomCity() enums.City {
	// fallback ثابت و امن (برای تست)
	fallback := []enums.City{
		enums.City(0),
		enums.City(1),
		enums.City(2),
		enums.City(3),
		enums.City(4),
	}
	return fallback[rand.Intn(len(fallback))]
}

func randomProvince() enums.Province {
	fallback := []enums.Province{
		enums.Province(0),
		enums.Province(1),
		enums.Province(2),
		enums.Province(3),
	}
	return fallback[rand.Intn(len(fallback))]
}
