package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"not null;index" json:"user_id"`
	ShiftID        uint           `gorm:"not null;index" json:"shift_id"`
	Date           time.Time      `gorm:"type:date;not null;index" json:"date"`
	OfficeID       uint           `gorm:"not null;index" json:"office_id"`
	Status         string         `gorm:"type:varchar(20);not null;default:'absent'" json:"status"`

	// Check-in
	TimeIn         *time.Time     `json:"time_in"`
	LatIn          *float64       `json:"lat_in"`
	LngIn          *float64       `json:"lng_in"`
	ImageIn        string         `gorm:"type:varchar(255)" json:"image_in"`
	DistanceMeters *float64       `json:"distance_meters"`

	// Check-out
	TimeOut        *time.Time     `json:"time_out"`
	LatOut         *float64       `json:"lat_out"`
	LngOut         *float64       `json:"lng_out"`
	ImageOut       string         `gorm:"type:varchar(255)" json:"image_out"`
	Duration       string         `gorm:"type:varchar(20)" json:"duration"`
	OvertimeMinutes int           `gorm:"not null;default:0" json:"overtime_minutes"`

	// Correction
	CorrectedBy    *uint          `gorm:"index" json:"corrected_by,omitempty"`
	CorrectedAt    *time.Time     `json:"corrected_at,omitempty"`
	CorrectionReason string       `gorm:"type:text" json:"correction_reason,omitempty"`

	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User   User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Shift  Shift         `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
	Office OfficeLocation `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
}

func (Attendance) TableName() string {
	return "attendances"
}

func (Attendance) GORMIndexes() []gorm.Index {
	return []gorm.Index{
		&AttendanceUserDateShiftIndex{},
	}
}

// AttendanceUserDateShiftIndex defines the unique constraint (user_id, date, shift_id).
type AttendanceUserDateShiftIndex struct {
	gorm.Index
}

func (AttendanceUserDateShiftIndex) TableName() string {
	return "attendances"
}

func (AttendanceUserDateShiftIndex) Name() string {
	return "idx_user_date_shift"
}

func (AttendanceUserDateShiftIndex) IsUnique() bool {
	return true
}

func (AttendanceUserDateShiftIndex) Columns() []string {
	return []string{"user_id", "date", "shift_id"}
}
