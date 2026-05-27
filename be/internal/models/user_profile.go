package models

import (
	"time"
)

type UserProfile struct {
	UserID        uint           `gorm:"primaryKey" json:"user_id"`
	Department    string         `gorm:"type:varchar(100)" json:"department"`
	Position      string         `gorm:"type:varchar(100)" json:"position"`
	Phone         string         `gorm:"type:varchar(20)" json:"phone"`
	JoinDate      *time.Time     `json:"join_date"`
	FacePhotoURL  string         `gorm:"type:varchar(500)" json:"face_photo_url"`
	FaceEmbedding string         `gorm:"type:text" json:"face_embedding"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
