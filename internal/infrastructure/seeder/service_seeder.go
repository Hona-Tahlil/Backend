package seeder

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
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

	for i := 0; i < 50; i++ {
		description := descriptions[rand.Intn(len(descriptions))]
		service := entities.Service{
			PetSitterID: uint(rand.Intn(20) + 1),             // Assuming 20 pet sitters exist
			Type:        enums.ServiceType(rand.Intn(2) + 1), // Assuming 5 service types
			Price:       uint(rand.Intn(90)+10) * 1000,       // Price between 10k-100k
			Description: &description,
		}

		if err := db.Create(&service).Error; err != nil {
			return err
		}
	}

	return nil
}
