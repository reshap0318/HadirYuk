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

// FindHistory finds attendance records with filters and pagination.
func (r *AttendanceRepository) FindHistory(tx *gorm.DB, userID uint, dateFrom, dateTo *time.Time, status string, page, pageSize int, preloads ...string) (*PagedResult[models.Attendance], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	db := r.getDB(tx)
	query := db.Model(&models.Attendance{}).Where("user_id = ?", userID)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	query = query.Order("date DESC, time_in DESC")

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	var attendances []models.Attendance
	if err := query.Limit(pageSize).Offset(offset).Find(&attendances).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[models.Attendance]{
		Data:       attendances,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// FindHistoryAll finds attendance records for all users with filters and pagination.
func (r *AttendanceRepository) FindHistoryAll(tx *gorm.DB, dateFrom, dateTo *time.Time, status string, page, pageSize int, preloads ...string) (*PagedResult[models.Attendance], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	db := r.getDB(tx)
	query := db.Model(&models.Attendance{})

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	query = query.Order("date DESC, time_in DESC")

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	var attendances []models.Attendance
	if err := query.Limit(pageSize).Offset(offset).Find(&attendances).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[models.Attendance]{
		Data:       attendances,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// FindMonthlyStats calculates monthly attendance statistics for a user.
func (r *AttendanceRepository) FindMonthlyStats(tx *gorm.DB, userID uint, year, month int) (present, late, absent, totalOvertime int, err error) {
	db := r.getDB(tx)
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	type Stats struct {
		Present  int `gorm:"column:present"`
		Late     int `gorm:"column:late"`
		Absent   int `gorm:"column:absent"`
		Overtime int `gorm:"column:overtime"`
	}

	var stats Stats
	err = db.Raw(`
		SELECT
			COALESCE(SUM(CASE WHEN status = 'present' THEN 1 ELSE 0 END), 0) as present,
			COALESCE(SUM(CASE WHEN status = 'late' THEN 1 ELSE 0 END), 0) as late,
			COALESCE(SUM(CASE WHEN status = 'absent' THEN 1 ELSE 0 END), 0) as absent,
			COALESCE(SUM(overtime_minutes), 0) as overtime
		FROM attendances
		WHERE user_id = ? AND date >= ? AND date < ? AND deleted_at IS NULL
	`, userID, startDate, endDate).Scan(&stats).Error

	return stats.Present, stats.Late, stats.Absent, stats.Overtime, err
}

// FindLateStats finds late attendance records with filters.
func (r *AttendanceRepository) FindLateStats(tx *gorm.DB, dateFrom, dateTo *time.Time, userID *uint, page, pageSize int, preloads ...string) (*PagedResult[models.Attendance], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	db := r.getDB(tx)
	query := db.Model(&models.Attendance{}).Where("status = ?", "late")

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	query = query.Order("date DESC, time_in DESC")

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	var attendances []models.Attendance
	if err := query.Limit(pageSize).Offset(offset).Find(&attendances).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[models.Attendance]{
		Data:       attendances,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
