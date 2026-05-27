package models

import (
	"time"

	"gorm.io/gorm"
)

type LeaveType struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description  *string        `gorm:"type:text" json:"description"`
	DefaultDays  int            `gorm:"not null;default:0" json:"default_days"`
	IsPaid       bool           `gorm:"not null;default:true" json:"is_paid"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LeaveType) TableName() string {
	return "leave_types"
}
