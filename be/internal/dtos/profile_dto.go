package dtos

import (
	"time"

	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
)

// ProfileUpdateRequest represents the request to update profile.
type ProfileUpdateRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Phone      string `json:"phone" validate:"omitempty"`
	Department string `json:"department" validate:"omitempty"`
	Position   string `json:"position" validate:"omitempty"`
	Avatar     string `json:"avatar"`
}

// FacePhotoRequest represents the request to upload face photo.
type FacePhotoRequest struct {
	FacePhoto string `json:"face_photo"`
}

// ProfileDTO represents profile data transfer object.
type ProfileDTO struct {
	ID         uint       `json:"id"`
	Email      string     `json:"email"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	Department string     `json:"department"`
	Position   string     `json:"position"`
	JoinDate   *time.Time `json:"join_date"`
	Avatar     string     `json:"avatar"`
	CreatedAt  time.Time  `json:"created_at"`
	Roles      []RoleMiniDTO   `json:"roles"`
	Permissions []PermissionDTO `json:"permissions"`
}

// ToProfileDTO converts User model to ProfileDTO.
func ToProfileDTO(u *models.User) ProfileDTO {
	dto := ProfileDTO{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Avatar:    helpers.GetFileURL(u.Avatar),
		CreatedAt: u.CreatedAt,
		Roles:     []RoleMiniDTO{},
		Permissions: []PermissionDTO{},
	}

	if u.Profile != nil {
		dto.Phone = u.Profile.Phone
		dto.Department = u.Profile.Department
		dto.Position = u.Profile.Position
		dto.JoinDate = u.Profile.JoinDate
	}

	permSet := make(map[uint]bool)
	for _, r := range u.Roles {
		dto.Roles = append(dto.Roles, ToRoleMiniDTO(&r))
		for _, p := range r.Permissions {
			if !permSet[p.ID] {
				permSet[p.ID] = true
				dto.Permissions = append(dto.Permissions, ToPermissionDTO(&p))
			}
		}
	}

	return dto
}
