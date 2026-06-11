package models

import (
	"time"

	"gorm.io/gorm"
)

type QRCode struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	OfficeID   uint           `gorm:"not null;index" json:"office_id"`
	CodeValue  string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"code_value"`
	Signature  string         `gorm:"type:varchar(255);not null" json:"signature"`
	ExpiresAt  time.Time      `gorm:"not null" json:"expires_at"`
	IsActive   bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedBy  uint           `gorm:"not null;index" json:"created_by"`
	RevokedAt  *time.Time     `json:"revoked_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Office OfficeLocation `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	Creator User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (QRCode) TableName() string {
	return "qr_codes"
}
