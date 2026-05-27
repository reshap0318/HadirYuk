package models

import (
	"time"

	"gorm.io/gorm"
)

type Shift struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	StartTime     string         `gorm:"type:varchar(5);not null" json:"start_time"`
	EndTime       string         `gorm:"type:varchar(5);not null" json:"end_time"`
	BreakDuration int            `gorm:"not null;default:0" json:"break_duration"`
	ColorCode     string         `gorm:"type:varchar(7);not null" json:"color_code"`
	TotalHours    float64        `gorm:"type:decimal(5,2)" json:"total_hours"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Shift) TableName() string {
	return "shifts"
}
