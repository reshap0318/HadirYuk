package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// SeedOfficeLocations inserts default office location data.
func SeedOfficeLocations(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding office locations...")

	locations := []struct {
		Name         string
		Address      string
		Latitude     float64
		Longitude    float64
		RadiusMeters int
		IsActive     bool
	}{
		{
			Name:         "Kantor Pusat",
			Address:      "Jl. Jend. Sudirman No. 1, Jakarta Pusat",
			Latitude:     -6.181788,
			Longitude:    106.820668,
			RadiusMeters: 120,
			IsActive:     true,
		},
	}

	resultMap := make(map[string]uint)

	for _, locData := range locations {
		var existing models.OfficeLocation
		err := db.Where("name = ?", locData.Name).First(&existing).Error
		if err == nil {
			resultMap[locData.Name] = existing.ID
			fmt.Printf("  ⊘ Office location %s already exists, skipping\n", locData.Name)
			continue
		}

		location := models.OfficeLocation{
			Name:         locData.Name,
			Address:      locData.Address,
			Latitude:     locData.Latitude,
			Longitude:    locData.Longitude,
			RadiusMeters: locData.RadiusMeters,
			IsActive:     locData.IsActive,
		}

		if err := db.Create(&location).Error; err != nil {
			log.Printf("Failed to create office location %s: %v", locData.Name, err)
		} else {
			resultMap[location.Name] = location.ID
		}
	}

	fmt.Printf("✓ Seeded %d office locations\n", len(resultMap))
	return resultMap
}
