package repositories

import (
	"time"

	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	*GenericRepository[models.Attendance]
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{
		GenericRepository: NewGenericRepository(db, &models.Attendance{}),
	}
}

// FindByUserAndDate finds attendance record for a specific user on a specific date.
func (r *AttendanceRepository) FindByUserAndDate(tx *gorm.DB, userID uint, date time.Time) (*models.Attendance, error) {
	db := r.getDB(tx)
	var attendance models.Attendance

	// Normalize date to start of day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := db.Model(&attendance).
		Where("user_id = ? AND date >= ? AND date < ?", userID, startOfDay, endOfDay).
		First(&attendance).Error

	if err != nil {
		return nil, err
	}

	return &attendance, nil
}
