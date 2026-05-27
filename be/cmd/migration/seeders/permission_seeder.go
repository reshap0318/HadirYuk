package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// SeedPermissions inserts default permission data.
func SeedPermissions(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding permissions...")

	permissions := []struct {
		Name        string
		Description string
	}{
		// User
		{"user.index", "View users list"},
		{"user.create", "Create new user"},
		{"user.edit", "Update user"},
		{"user.delete", "Delete user"},
		// Role
		{"role.index", "View roles list"},
		{"role.create", "Create new role"},
		{"role.edit", "Update role"},
		{"role.delete", "Delete role"},
		// Permission
		{"permission.index", "View permissions list"},
		{"permission.create", "Create new permission"},
		{"permission.edit", "Update permission"},
		{"permission.delete", "Delete permission"},
	}

	resultMap := make(map[string]uint)

	for _, perm := range permissions {
		var existing models.Permission
		err := db.Where("name = ?", perm.Name).First(&existing).Error
		if err == nil {
			resultMap[perm.Name] = existing.ID
			fmt.Printf("  ⊘ Permission %s already exists, skipping\n", perm.Name)
			continue
		}

		p := models.Permission{
			Name:        perm.Name,
			Description: strPtr(perm.Description),
		}

		if err := db.Create(&p).Error; err != nil {
			log.Printf("Failed to create permission %s: %v", perm.Name, err)
		} else {
			resultMap[perm.Name] = p.ID
		}
	}

	fmt.Printf("✓ Seeded %d permissions\n", len(resultMap))
	return resultMap
}

func strPtr(s string) *string {
	return &s
}
