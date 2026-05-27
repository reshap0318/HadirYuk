package seeders

import (
	"fmt"
	"log"
	"time"

	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// SeedUsers inserts default user data.
func SeedUsers(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding users...")

	defaultUsers := []struct {
		Email      string
		Password   string
		Name       string
		Phone      string
		Department string
		Position   string
		JoinDate   time.Time
	}{
		{
			Email:      "suAdmin@app.com",
			Password:   "@dmin#123",
			Name:       "Super Administrator",
			Phone:      "081234567890",
			Department: "IT",
			Position:   "System Administrator",
			JoinDate:   time.Now(),
		},
		{
			Email:      "hradmin@app.com",
			Password:   "HrAdmin#123",
			Name:       "HR Administrator",
			Phone:      "081234567891",
			Department: "Human Resources",
			Position:   "HR Manager",
			JoinDate:   time.Now(),
		},
		{
			Email:      "karyawan@app.com",
			Password:   "Karyawan#123",
			Name:       "Karyawan Demo",
			Phone:      "081234567892",
			Department: "Engineering",
			Position:   "Software Engineer",
			JoinDate:   time.Now(),
		},
	}

	resultMap := make(map[string]uint)

	for _, userData := range defaultUsers {
		var existing models.User
		result := db.Where("email = ?", userData.Email).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			hashedPassword, err := helpers.HashString(userData.Password)
			if err != nil {
				log.Printf("Failed to hash password for %s: %v", userData.Email, err)
				continue
			}

			user := models.User{
				Email:    userData.Email,
				Password: hashedPassword,
				Name:     userData.Name,
				Profile: &models.UserProfile{
					Phone:      userData.Phone,
					Department: userData.Department,
					Position:   userData.Position,
					JoinDate:   &userData.JoinDate,
				},
			}

			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user %s: %v", userData.Email, err)
			} else {
				resultMap[user.Email] = user.ID
			}
		} else if result.Error != nil {
			log.Printf("Failed to check user %s: %v", userData.Email, result.Error)
		} else {
			resultMap[userData.Email] = existing.ID
			fmt.Printf("  ⊘ User %s already exists, skipping\n", userData.Email)

			var profile models.UserProfile
			if db.Where("user_id = ?", existing.ID).First(&profile).Error == gorm.ErrRecordNotFound {
				profile = models.UserProfile{
					UserID:     existing.ID,
					Phone:      userData.Phone,
					Department: userData.Department,
					Position:   userData.Position,
					JoinDate:   &userData.JoinDate,
				}
				if err := db.Create(&profile).Error; err != nil {
					log.Printf("Failed to create profile for %s: %v", userData.Email, err)
				} else {
					fmt.Printf("    + Created profile for %s\n", userData.Email)
				}
			}
		}
	}

	fmt.Printf("✓ Seeded %d users\n", len(resultMap))
	return resultMap
}
