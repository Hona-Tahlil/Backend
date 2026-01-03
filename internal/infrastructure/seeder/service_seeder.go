package seeder

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"log/slog"
	"math/rand"

	"gorm.io/gorm"
)

func ServiceSeeder(db *gorm.DB) error {
	descriptions := []string{
		"Basic pet sitting service",
		"Premium pet care with walks",
		"Overnight pet sitting",
		"Pet grooming service",
		"Pet training session",
	}
	if db.First(&entities.Service{}).Error == nil {
		slog.Info("✓ Services already seeded, skipping...")
		return nil
	}

	for i := 0; i < 50; i++ {
		description := descriptions[rand.Intn(len(descriptions))]
		service := entities.Service{
			PetSitterID: uint(rand.Intn(20) + 1),
			Type:        enums.ServiceType(rand.Intn(2) + 1),
			Price:       uint(rand.Intn(90)+10) * 1000,
			Description: &description,
		}

		if err := db.Create(&service).Error; err != nil {
			return err
		}
	}

	slog.Info("✓ Successfully seeded services", "count", 50)

	return nil
}
