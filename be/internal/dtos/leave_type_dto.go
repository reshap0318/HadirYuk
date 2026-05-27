package dtos

import "github.com/reshap0318/hadirYuk/internal/models"

type LeaveTypeRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=50"`
	Description *string `json:"description"`
	DefaultDays int     `json:"default_days" validate:"required,min=0"`
	IsPaid      bool    `json:"is_paid"`
}

type LeaveTypeDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	DefaultDays int     `json:"default_days"`
	IsPaid      bool    `json:"is_paid"`
}

func ToLeaveTypeDTO(l *models.LeaveType) LeaveTypeDTO {
	return LeaveTypeDTO{
		ID:          l.ID,
		Name:        l.Name,
		Description: l.Description,
		DefaultDays: l.DefaultDays,
		IsPaid:      l.IsPaid,
	}
}

func ToLeaveTypeDTOList(leaveTypes []models.LeaveType) []LeaveTypeDTO {
	result := make([]LeaveTypeDTO, len(leaveTypes))
	for i, l := range leaveTypes {
		result[i] = ToLeaveTypeDTO(&l)
	}
	return result
}
