package models

import (
	"time"

	"gorm.io/gorm"
)

type UserShiftAssignment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	ShiftID   uint           `gorm:"not null;index" json:"shift_id"`
	StartDate time.Time      `gorm:"not null" json:"start_date"`
	EndDate   *time.Time     `json:"end_date"` // nil = ongoing
	IsActive  bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Shift Shift `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
}

func (UserShiftAssignment) TableName() string {
	return "user_shift_assignments"
}
