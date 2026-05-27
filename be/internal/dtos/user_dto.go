package dtos

import (
	"time"

	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
)

// UserCreateRequest represents the request to create a user.
type UserCreateRequest struct {
	Email            string     `json:"email" validate:"required,email"`
	Password         string     `json:"password" validate:"required,min=6"`
	Name             string     `json:"name" validate:"required,min=2,max=100"`
	Phone            string     `json:"phone" validate:"omitempty"`
	Department       string     `json:"department" validate:"omitempty"`
	Position         string     `json:"position" validate:"omitempty"`
	JoinDate         *time.Time `json:"join_date" validate:"omitempty"`
	Avatar           string     `json:"avatar"`
	FacePhoto        string     `json:"face_photo"`
	Roles            []uint     `json:"roles"`
}

// UserUpdateRequest represents the request to update a user.
type UserUpdateRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Phone      string `json:"phone" validate:"omitempty"`
	Department string `json:"department" validate:"omitempty"`
	Position   string `json:"position" validate:"omitempty"`
	Avatar     string `json:"avatar"`
	Roles      []uint `json:"roles"`
}

// UserStatusUpdateRequest represents the request to update user status.
type UserStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended"`
}

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

// UserDTO represents user data transfer object.
type UserDTO struct {
	ID          uint            `json:"id"`
	Email       string          `json:"email"`
	Name        string          `json:"name"`
	Phone       string          `json:"phone"`
	Department  string          `json:"department"`
	Position    string          `json:"position"`
	JoinDate    *time.Time      `json:"join_date"`
	Avatar      string          `json:"avatar"`
	FacePhoto   string          `json:"face_photo_url"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Roles       []RoleMiniDTO   `json:"roles"`
	Permissions []PermissionDTO `json:"permissions"`
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

// ToUserDTO converts User model to UserDTO.
func ToUserDTO(u *models.User) UserDTO {
	dto := UserDTO{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Avatar:      helpers.GetFileURL(u.Avatar),
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Roles:       []RoleMiniDTO{},
		Permissions: []PermissionDTO{},
	}

	if u.Profile != nil {
		dto.Phone = u.Profile.Phone
		dto.Department = u.Profile.Department
		dto.Position = u.Profile.Position
		dto.JoinDate = u.Profile.JoinDate
		dto.FacePhoto = helpers.GetFileURL(u.Profile.FacePhotoURL)
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

// ToUserDTOList converts a slice of User models to UserDTOs.
func ToUserDTOList(users []models.User) []UserDTO {
	result := make([]UserDTO, len(users))
	for i, u := range users {
		result[i] = ToUserDTO(&u)
	}
	return result
}
