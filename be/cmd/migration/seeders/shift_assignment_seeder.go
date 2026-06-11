package seeders

import (
	"fmt"
	"log"
	"time"

	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// SeedShiftAssignments assigns the "Pagi" shift to all users.
func SeedShiftAssignments(db *gorm.DB, shiftIDs map[string]uint) {
	fmt.Println("Seeding shift assignments...")

	shiftID, shiftOK := shiftIDs["Pagi"]
	if !shiftOK {
		log.Println("Shift 'Pagi' not found, skipping shift assignments")
		return
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Failed to fetch users: %v", err)
		return
	}

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(1, 0, 0)

	count := 0
	for _, user := range users {
		var existing models.UserShiftAssignment
		err := db.Where("user_id = ? AND shift_id = ? AND start_date = ?", user.ID, shiftID, startDate).First(&existing).Error
		if err == nil {
			fmt.Printf("  ⊘ Shift assignment for user %d already exists, skipping\n", user.ID)
			continue
		}

		assignment := models.UserShiftAssignment{
			UserID:    user.ID,
			ShiftID:   shiftID,
			StartDate: startDate,
			EndDate:   &endDate,
			IsActive:  true,
		}

		if err := db.Create(&assignment).Error; err != nil {
			log.Printf("Failed to create shift assignment for user %d: %v", user.ID, err)
		} else {
			count++
		}
	}

	fmt.Printf("✓ Seeded %d shift assignments\n", count)
}
