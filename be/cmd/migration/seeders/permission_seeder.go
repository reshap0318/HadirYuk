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
		// Auth
		{"auth.login", "Can login to system"},
		{"auth.logout", "Can logout from system"},
		{"auth.forgot-password", "Can request password reset"},
		{"auth.reset-password", "Can reset password with token"},
		{"auth.change-password", "Can change own password"},
		// Profile
		{"profile.view", "Can view own profile"},
		{"profile.update", "Can update own profile"},
		{"profile.upload-face", "Can upload face photo"},
		// Attendance
		{"attendance.checkin", "Can check-in attendance"},
		{"attendance.checkout", "Can check-out attendance"},
		{"attendance.view", "Can view own attendance history"},
		{"attendance.view-all", "Can view all attendance history"},
		{"attendance.export", "Can export attendance report"},
		{"attendance.correct", "Can correct attendance record"},
		// Shift
		{"shift.index", "Can view shift list"},
		{"shift.create", "Can create new shift"},
		{"shift.update", "Can update shift"},
		{"shift.delete", "Can delete shift"},
		{"shift.assign", "Can assign shift to employee"},
		// Leave
		{"leave.submit", "Can submit leave request"},
		{"leave.view", "Can view own leave history"},
		{"leave.view-all", "Can view all leave history"},
		{"leave.manage-types", "Can manage leave types"},
		// User
		{"user.index", "Can view user list"},
		{"user.create", "Can create new user"},
		{"user.update", "Can update user"},
		{"user.delete", "Can delete/deactivate user"},
		{"user.assign-role", "Can assign role to user"},
		{"user.view-all", "Can view all user data including super admin"},
		// Role
		{"role.index", "Can view role list"},
		{"role.create", "Can create new role"},
		{"role.update", "Can update role"},
		{"role.delete", "Can delete role"},
		{"role.assign-permission", "Can assign permission to role"},
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
		{"dashboard.view", "Can view own dashboard"},
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
