package seeder

import (
	"hona/backend/internal/domain/entities"
	"log"

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

	if err := s.SeedUsers(50); err != nil {
		return err
	}

	if err := s.SeedPets(150); err != nil {
		return err
	}

	if err := s.SeedPetSitters(20); err != nil {
		return err
	}

	if err := s.SeedServices(30); err != nil {
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
	// TODO: implement when needed
	log.Println("⚠️  Wallet seeding not implemented yet")
	return nil
}

// SeedAddresses seeds address data
func (s *DatabaseSeeder) SeedAddresses(count int) error {
	log.Printf("📍 Seeding %d addresses...", count)
	// TODO: implement when needed
	log.Println("⚠️  Address seeding not implemented yet")
	return nil
}

// SeedChats seeds chat data
func (s *DatabaseSeeder) SeedChats(count int) error {
	log.Printf("💬 Seeding %d chats...", count)
	// TODO: implement when needed
	log.Println("⚠️  Chat seeding not implemented yet")
	return nil
}

func (s *DatabaseSeeder) ClearAll() {
	log.Println("🗑️  Clearing all tables...")

	s.db.Migrator().DropTable(
		&entities.User{},
		&entities.Role{},
		&entities.Permission{},
		&entities.Wallet{},
		&entities.Request{},
		&entities.CalendarSlot{},
		&entities.Pet{},
		&entities.PetSitter{},
		&entities.Service{},
		&entities.Chat{},
		&entities.Comment{},
		&entities.Address{},
		&entities.TextMessage{},
		&entities.Transaction{},
		&entities.Transfer{},
	)

}
