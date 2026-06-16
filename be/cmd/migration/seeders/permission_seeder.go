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
		// Attendance
		{"attendance.checkin", "Can check-in (geotagging)"},
		{"attendance.checkout", "Can check-out (geotagging)"},
		{"attendance.view-all", "Can view all attendance history (NOT YET IMPLEMENTED)"},
		{"attendance.export", "Can export attendance report (NOT YET IMPLEMENTED)"},
		{"attendance.correct", "Can correct attendance record (NOT YET IMPLEMENTED)"},
		// Shift
		{"shift.index", "Can view shift list"},
		{"shift.create", "Can create new shift"},
		{"shift.update", "Can update shift"},
		{"shift.delete", "Can delete shift"},
		{"shift-assign.index", "Can view shift assignments"},
		{"shift-assign.create", "Can assign shift to employee"},
		{"shift-assign.update", "Can update shift assignment"},
		{"shift-assign.delete", "Can delete shift assignment"},
		// User
		{"user.index", "Can view user list"},
		{"user.view-all", "Can view all user data including super admin"},
		{"user.create", "Can create new user"},
		{"user.update", "Can update user"},
		{"user.delete", "Can delete/deactivate user"},
		// Role
		{"role.index", "Can view role list"},
		{"role.index-all", "Can view all roles including Super Admin"},
		{"role.create", "Can create new role"},
		{"role.update", "Can update role"},
		{"role.delete", "Can delete role"},
		// Permission
		{"permission.index", "Can view permission list"},
		{"permission.create", "Can create new permission"},
		{"permission.update", "Can update permission"},
		{"permission.delete", "Can delete permission"},
		// Location
		{"location.index", "Can view location list"},
		{"location.create", "Can create new location"},
		{"location.update", "Can update location"},
		{"location.delete", "Can delete location"},
		// Dashboard
		{"dashboard.view-hr", "Can view HR dashboard"},
		{"dashboard.view-admin", "Can view admin dashboard"},
		// Report
		{"report.view", "Can view reports"},
		{"report.export-excel", "Can export report to Excel"},
		{"report.export-pdf", "Can export report to PDF"},
		// QR Code
		{"qrcode.generate", "Can generate QR code"},
		{"qrcode.view", "Can view active QR codes"},
		{"qrcode.revoke", "Can revoke QR code"},
		// Audit & Statistics
		{"audit.view", "Can view audit log"},
		{"late-statistic.view", "Can view late statistics"},
	}

	resultMap := make(map[string]uint)

	for _, perm := range permissions {
		var existing models.Permission
		err := db.Where("name = ?", perm.Name).First(&existing).Error
		if err == nil {
			resultMap[perm.Name] = existing.ID
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
