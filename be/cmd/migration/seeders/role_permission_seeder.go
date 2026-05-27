package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// SeedRolePermissions maps roles to their permissions.
func SeedRolePermissions(db *gorm.DB, roleIDs map[string]uint, permIDs map[string]uint) {
	fmt.Println("Seeding role permissions...")

	rolePerms := map[string][]string{
		"Super Admin": {},
		"HR Admin": {
			// Auth & Profile
			"auth.login", "auth.logout", "auth.forgot-password", "auth.reset-password", "auth.change-password",
			"profile.view", "profile.update", "profile.upload-face",
			// Attendance
			"attendance.checkin", "attendance.checkout", "attendance.view", "attendance.view-all",
			"attendance.export", "attendance.correct",
			// Shift
			"shift.index", "shift.create", "shift.update", "shift.delete", "shift.assign",
			// Leave
			"leave.submit", "leave.view", "leave.view-all", "leave.manage-types",
			// User
			"user.index", "user.create", "user.update",
			// Location
			"location.index", "location.create", "location.update", "location.delete",
			// Dashboard
			"dashboard.view", "dashboard.view-hr",
			// Report
			"report.view", "report.export-excel", "report.export-pdf",
			// QR Code
			"qrcode.generate", "qrcode.view", "qrcode.revoke",
			// Statistics
			"late-statistic.view",
		},
		"Karyawan": {
			// Auth & Profile
			"auth.login", "auth.logout", "auth.forgot-password", "auth.reset-password", "auth.change-password",
			"profile.view", "profile.update", "profile.upload-face",
			// Attendance
			"attendance.checkin", "attendance.checkout", "attendance.view",
			// Shift
			"shift.index",
			// Leave
			"leave.submit", "leave.view",
			// Dashboard
			"dashboard.view",
		},
	}

	count := 0
	for roleName, permNames := range rolePerms {
		roleID, roleOK := roleIDs[roleName]
		if !roleOK {
			log.Printf("Role %s not found, skipping", roleName)
			continue
		}

		// Super Admin gets ALL permissions
		if roleName == "Super Admin" {
			for permName := range permIDs {
				permID := permIDs[permName]

				var existing models.RoleHasPermission
				err := db.Where("role_id = ? AND permission_id = ?", roleID, permID).First(&existing).Error
				if err == nil {
					continue
				}

				rp := models.RoleHasPermission{
					RoleID:       roleID,
					PermissionID: permID,
				}

				if err := db.Create(&rp).Error; err != nil {
					log.Printf("Failed to create role_permission for %s-%s: %v", roleName, permName, err)
				} else {
					count++
				}
			}
			continue
		}

		for _, permName := range permNames {
			permID, permOK := permIDs[permName]
			if !permOK {
				log.Printf("Permission %s not found, skipping", permName)
				continue
			}

			var existing models.RoleHasPermission
			err := db.Where("role_id = ? AND permission_id = ?", roleID, permID).First(&existing).Error
			if err == nil {
				continue
			}

			rp := models.RoleHasPermission{
				RoleID:       roleID,
				PermissionID: permID,
			}

			if err := db.Create(&rp).Error; err != nil {
				log.Printf("Failed to create role_permission for %s-%s: %v", roleName, permName, err)
			} else {
				count++
			}
		}
	}

	fmt.Printf("✓ Seeded %d role permissions\n", count)
}
