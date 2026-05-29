package dtos

import (
	"time"

	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
)

type ShiftAssignmentRequest struct {
	UserID    uint   `json:"user_id" validate:"required"`
	ShiftID   uint   `json:"shift_id" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date"`
}

type ShiftAssignmentUpdateRequest struct {
	ShiftID   uint   `json:"shift_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsActive  *bool  `json:"is_active"`
}

type ShiftAssignmentDTO struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	ShiftID   uint               `json:"shift_id"`
	StartDate time.Time          `json:"start_date"`
	EndDate   *time.Time         `json:"end_date"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	User      *ShiftAssignUserInfo  `json:"user,omitempty"`
	Shift     *ShiftDTO          `json:"shift,omitempty"`
}

type ShiftAssignUserInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

func ToShiftAssignmentDTO(a *models.UserShiftAssignment) ShiftAssignmentDTO {
	dto := ShiftAssignmentDTO{
		ID:        a.ID,
		UserID:    a.UserID,
		ShiftID:   a.ShiftID,
		StartDate: a.StartDate,
		EndDate:   a.EndDate,
		IsActive:  a.IsActive,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	if a.User.ID != 0 {
		dto.User = &ShiftAssignUserInfo{
			ID:     a.User.ID,
			Name:   a.User.Name,
			Email:  a.User.Email,
			Avatar: helpers.GetFileURL(a.User.Avatar),
		}
	}

	if a.Shift.ID != 0 {
		shiftDTO := ToShiftDTO(&a.Shift)
		dto.Shift = &shiftDTO
	}

	return dto
}

func ToShiftAssignmentDTOList(assignments []models.UserShiftAssignment) []ShiftAssignmentDTO {
	result := make([]ShiftAssignmentDTO, len(assignments))
	for i, a := range assignments {
		result[i] = ToShiftAssignmentDTO(&a)
	}
	return result
}
