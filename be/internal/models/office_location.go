package models

import (
	"time"

	"gorm.io/gorm"
)

type OfficeLocation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null" json:"name"`
	Address       string         `gorm:"type:text;not null" json:"address"`
	Latitude      float64        `gorm:"not null" json:"latitude"`
	Longitude     float64        `gorm:"not null" json:"longitude"`
	RadiusMeters  int            `gorm:"not null" json:"radius_meters"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (OfficeLocation) TableName() string {
	return "office_locations"
}
