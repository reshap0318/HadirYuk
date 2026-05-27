package dtos

import "github.com/reshap0318/hadirYuk/internal/models"

type ShiftRequest struct {
	Name          string  `json:"name" validate:"required,min=3,max=50"`
	StartTime     string  `json:"start_time" validate:"required"`
	EndTime       string  `json:"end_time" validate:"required"`
	BreakDuration int     `json:"break_duration" validate:"required,min=0"`
	ColorCode     string  `json:"color_code" validate:"required"`
	TotalHours    float64 `json:"total_hours"`
}

type ShiftDTO struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	BreakDuration int     `json:"break_duration"`
	ColorCode     string  `json:"color_code"`
	TotalHours    float64 `json:"total_hours"`
}

func ToShiftDTO(s *models.Shift) ShiftDTO {
	return ShiftDTO{
		ID:            s.ID,
		Name:          s.Name,
		StartTime:     s.StartTime,
		EndTime:       s.EndTime,
		BreakDuration: s.BreakDuration,
		ColorCode:     s.ColorCode,
		TotalHours:    s.TotalHours,
	}
}

func ToShiftDTOList(shifts []models.Shift) []ShiftDTO {
	result := make([]ShiftDTO, len(shifts))
	for i, s := range shifts {
		result[i] = ToShiftDTO(&s)
	}
	return result
}
