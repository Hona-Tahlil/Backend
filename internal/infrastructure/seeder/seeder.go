package seeder

import (
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

	// Seed in order of dependencies
	if err := s.SeedUsers(50); err != nil {
		return err
	}

	if err := s.SeedPets(150); err != nil {
		return err
	}

	// Add more seeders as needed...
	// if err := s.SeedPetSitters(20); err != nil {
	// 	return err
	// }

	// if err := s.SeedServices(30); err != nil {
	// 	return err
	// }

	// if err := s.SeedRequests(100); err != nil {
	// 	return err
	// }

	log.Println("✅ Database seeding completed!")
	return nil
}

// SeedUsers seeds user data
func (s *DatabaseSeeder) SeedUsers(count int) error {
	log.Printf("👥 Seeding %d users...", count)
	seeder := NewUserSeeder(s.db)
	return seeder.Seed(count)
}

// SeedPets seeds pet data
func (s *DatabaseSeeder) SeedPets(count int) error {
	log.Printf("🐾 Seeding %d pets...", count)
	seeder := NewPetSeeder(s.db)
	return seeder.Seed(count)
}

// SeedPetSitters seeds pet sitter data
func (s *DatabaseSeeder) SeedPetSitters(count int) error {
	log.Printf("🏠 Seeding %d pet sitters...", count)
	// TODO: implement when needed
	log.Println("⚠️  Pet sitter seeding not implemented yet")
	return nil
}

// SeedServices seeds service data
func (s *DatabaseSeeder) SeedServices(count int) error {
	log.Printf("💼 Seeding %d services...", count)
	// TODO: implement when needed
	log.Println("⚠️  Service seeding not implemented yet")
	return nil
}

// SeedRequests seeds request data
func (s *DatabaseSeeder) SeedRequests(count int) error {
	log.Printf("📋 Seeding %d requests...", count)
	// TODO: implement when needed
	log.Println("⚠️  Request seeding not implemented yet")
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

// ClearAll truncates all tables (use with caution!)
func (s *DatabaseSeeder) ClearAll() error {
	log.Println("🗑️  Clearing all tables...")

	// Disable foreign key checks (PostgreSQL syntax)
	s.db.Exec("SET session_replication_role = 'replica'")

	// List all tables in dependency order (children first, parents last)
	tables := []string{
		// Children tables first
		"text_messages",
		"chats",
		"comments",
		"transfers",
		"transactions",
		"wallets",
		"request_services",
		"request_pets",
		"request_addresses",
		"requests",
		"calendar_slots",
		"pets",
		"pet_sitters",
		"addresses",
		"services",
		// Parent tables last
		"users",
		"cities",
		"provinces",
	}

	for _, table := range tables {
		// Use TRUNCATE CASCADE for PostgreSQL
		if err := s.db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			log.Printf("⚠️  Warning: failed to truncate %s: %v", table, err)
		} else {
			log.Printf("  ✓ Cleared table: %s", table)
		}
	}

	// Re-enable foreign key checks
	s.db.Exec("SET session_replication_role = 'origin'")

	log.Println("✓ All tables cleared")
	return nil
}

// SeedCustom allows custom seeding with specific counts for each entity
func (s *DatabaseSeeder) SeedCustom(userCount, petCount, petSitterCount, serviceCount, requestCount int) error {
	log.Println("🌱 Starting custom database seeding...")

	if userCount > 0 {
		if err := s.SeedUsers(userCount); err != nil {
			return err
		}
	}

	if petCount > 0 {
		if err := s.SeedPets(petCount); err != nil {
			return err
		}
	}

	if petSitterCount > 0 {
		if err := s.SeedPetSitters(petSitterCount); err != nil {
			return err
		}
	}

	if serviceCount > 0 {
		if err := s.SeedServices(serviceCount); err != nil {
			return err
		}
	}

	if requestCount > 0 {
		if err := s.SeedRequests(requestCount); err != nil {
			return err
		}
	}

	log.Println("✅ Custom seeding completed!")
	return nil
}

// SeedMinimal seeds minimal data for basic testing (faster)
func (s *DatabaseSeeder) SeedMinimal() error {
	log.Println("🌱 Starting minimal database seeding...")

	if err := s.SeedUsers(10); err != nil {
		return err
	}

	if err := s.SeedPets(20); err != nil {
		return err
	}

	log.Println("✅ Minimal seeding completed!")
	return nil
}

// SeedProduction seeds production-like data (more realistic volumes)
func (s *DatabaseSeeder) SeedProduction() error {
	log.Println("🌱 Starting production-like database seeding...")

	if err := s.SeedUsers(500); err != nil {
		return err
	}

	if err := s.SeedPets(1000); err != nil {
		return err
	}

	if err := s.SeedPetSitters(100); err != nil {
		return err
	}

	if err := s.SeedServices(50); err != nil {
		return err
	}

	if err := s.SeedRequests(300); err != nil {
		return err
	}

	log.Println("✅ Production-like seeding completed!")
	return nil
}

// GetStats returns seeding statistics
func (s *DatabaseSeeder) GetStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	tables := []string{
		"users",
		"pets",
		"pet_sitters",
		"services",
		"requests",
		"addresses",
		"wallets",
		"chats",
	}

	for _, table := range tables {
		var count int64
		if err := s.db.Table(table).Count(&count).Error; err != nil {
			return nil, err
		}
		stats[table] = count
	}

	return stats, nil
}

// PrintStats prints seeding statistics
func (s *DatabaseSeeder) PrintStats() error {
	stats, err := s.GetStats()
	if err != nil {
		return err
	}

	log.Println("\n📊 Database Statistics:")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for table, count := range stats {
		log.Printf("  %-20s: %d", table, count)
	}
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	return nil
}
