package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// SeedShifts inserts default shift data.
func SeedShifts(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding shifts...")

	shifts := []struct {
		Name         string
		StartTime    string
		EndTime      string
		BreakDuration int
		FlexiMinutes int
		ColorCode    string
		TotalHours   float64
	}{
		{"Pagi", "08:00", "17:00", 60, 15, "#3B82F6", 8.0},
		{"Siang", "14:00", "22:00", 60, 15, "#F59E0B", 8.0},
		{"Malam", "22:00", "06:00", 60, 15, "#8B5CF6", 8.0},
	}

	resultMap := make(map[string]uint)

	for _, shiftData := range shifts {
		var existing models.Shift
		err := db.Where("name = ?", shiftData.Name).First(&existing).Error
		if err == nil {
			resultMap[shiftData.Name] = existing.ID
			fmt.Printf("  ⊘ Shift %s already exists, skipping\n", shiftData.Name)
			continue
		}

		shift := models.Shift{
			Name:          shiftData.Name,
			StartTime:     shiftData.StartTime,
			EndTime:       shiftData.EndTime,
			BreakDuration: shiftData.BreakDuration,
			FlexiMinutes:  shiftData.FlexiMinutes,
			ColorCode:     shiftData.ColorCode,
			TotalHours:    shiftData.TotalHours,
		}

		if err := db.Create(&shift).Error; err != nil {
			log.Printf("Failed to create shift %s: %v", shiftData.Name, err)
		} else {
			resultMap[shift.Name] = shift.ID
		}
	}

	fmt.Printf("✓ Seeded %d shifts\n", len(resultMap))
	return resultMap
}
