package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// SeedRoles inserts default role data.
func SeedRoles(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding roles...")

	roles := []struct {
		Name        string
		Description string
	}{
		{"Super Admin", "Administrator tertinggi sistem dengan akses penuh ke semua fitur"},
		{"HR Admin", "Administrator operasional HR: karyawan, shift, absensi, cuti, laporan"},
		{"Karyawan", "Pengguna biasa untuk absensi, cuti, dan riwayat kehadiran pribadi"},
	}

	resultMap := make(map[string]uint)

	for _, roleData := range roles {
		var existing models.Role
		err := db.Where("name = ?", roleData.Name).First(&existing).Error
		if err == nil {
			resultMap[roleData.Name] = existing.ID
			fmt.Printf("  ⊘ Role %s already exists, skipping\n", roleData.Name)
			continue
		}

		role := models.Role{
			Name:        roleData.Name,
			Description: strPtr(roleData.Description),
		}

		if err := db.Create(&role).Error; err != nil {
			log.Printf("Failed to create role %s: %v", roleData.Name, err)
		} else {
			resultMap[roleData.Name] = role.ID
		}
	}

	fmt.Printf("✓ Seeded %d roles\n", len(resultMap))
	return resultMap
}
