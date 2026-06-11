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

// FindByUserAndDate finds attendance records for a specific user on a specific date.
func (r *AttendanceRepository) FindByUserAndDate(tx *gorm.DB, userID uint, date time.Time, preloads ...string) ([]models.Attendance, error) {
	db := r.getDB(tx)
	var attendances []models.Attendance

	// Normalize date to start of day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := db.Model(&models.Attendance{}).
		Where("user_id = ? AND date >= ? AND date < ?", userID, startOfDay, endOfDay)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.Find(&attendances).Error
	if err != nil {
		return nil, err
	}

	return attendances, nil
}

// FindActiveSessionByUserID finds the active session (TimeIn != nil AND TimeOut == nil) for a user across ALL dates.
func (r *AttendanceRepository) FindActiveSessionByUserID(tx *gorm.DB, userID uint, preloads ...string) (*models.Attendance, error) {
	db := r.getDB(tx)
	var attendance models.Attendance

	query := db.Model(&attendance).
		Where("user_id = ? AND time_in IS NOT NULL AND time_out IS NULL", userID)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&attendance).Error
	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

// FindByUserDateShift finds attendance record for a specific user, date, and shift.
func (r *AttendanceRepository) FindByUserDateShift(tx *gorm.DB, userID uint, date time.Time, shiftID uint, preloads ...string) (*models.Attendance, error) {
	db := r.getDB(tx)
	var attendance models.Attendance

	// Normalize date to start of day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := db.Model(&attendance).
		Where("user_id = ? AND date >= ? AND date < ? AND shift_id = ?", userID, startOfDay, endOfDay, shiftID)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&attendance).Error
	if err != nil {
		return nil, err
	}

	return &attendance, nil
}
