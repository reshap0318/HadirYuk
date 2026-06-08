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
	TimeIn         *time.Time     `json:"time_in"`
	Lat            *float64       `json:"lat"`
	Lng            *float64       `json:"lng"`
	OfficeID       uint           `gorm:"not null;index" json:"office_id"`
	Status         string         `gorm:"type:varchar(20);not null;default:'absent'" json:"status"`
	DistanceMeters *float64       `json:"distance_meters"`
	ImageIn        string         `gorm:"type:varchar(255)" json:"image_in"`
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
