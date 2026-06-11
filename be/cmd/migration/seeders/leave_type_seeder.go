package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// SeedLeaveTypes inserts default leave type data.
func SeedLeaveTypes(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding leave types...")

	leaveTypes := []struct {
		Name        string
		Description string
		DefaultDays int
		IsPaid      bool
	}{
		{"Cuti Tahunan", "Cuti tahunan yang diberikan kepada karyawan", 12, true},
		{"Sakit", "Cuti karena sakit dengan surat dokter", 0, true},
		{"Cuti Khusus", "Cuti untuk keperluan khusus (menikah, kelahiran, dll)", 3, true},
	}

	resultMap := make(map[string]uint)

	for _, ltData := range leaveTypes {
		var existing models.LeaveType
		err := db.Where("name = ?", ltData.Name).First(&existing).Error
		if err == nil {
			resultMap[ltData.Name] = existing.ID
			fmt.Printf("  ⊘ Leave type %s already exists, skipping\n", ltData.Name)
			continue
		}

		desc := ltData.Description
		leaveType := models.LeaveType{
			Name:        ltData.Name,
			Description: &desc,
			DefaultDays: ltData.DefaultDays,
			IsPaid:      ltData.IsPaid,
		}

		if err := db.Create(&leaveType).Error; err != nil {
			log.Printf("Failed to create leave type %s: %v", ltData.Name, err)
		} else {
			resultMap[leaveType.Name] = leaveType.ID
		}
	}

	fmt.Printf("✓ Seeded %d leave types\n", len(resultMap))
	return resultMap
}
