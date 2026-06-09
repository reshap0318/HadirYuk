package dtos

import (
	"time"

	"github.com/reshap0318/hadirYuk/internal/models"
)

type AttendanceCheckInRequest struct {
	Lat   float64 `json:"lat" validate:"required"`
	Lng   float64 `json:"lng" validate:"required"`
	Image string  `json:"image" validate:"required"`
}

type NearestOfficeRequest struct {
	Lat float64 `json:"latitude" validate:"required"`
	Lng float64 `json:"longitude" validate:"required"`
}

type NearestOfficeResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RadiusMeters int     `json:"radius_meters"`
	Distance     float64 `json:"distance"`
}

type AttendanceStatusResponse struct {
	HasCheckedIn  bool       `json:"has_checked_in"`
	TimeIn        *time.Time `json:"time_in,omitempty"`
	HasCheckedOut bool       `json:"has_checked_out"`
	TimeOut       *time.Time `json:"time_out,omitempty"`
	Duration      string     `json:"duration,omitempty"`
	ShiftID       uint       `json:"shift_id,omitempty"`
	Status        string     `json:"status,omitempty"`
	Distance      *float64   `json:"distance_meters,omitempty"`
}

type AttendanceDTO struct {
	ID             uint       `json:"id"`
	UserID         uint       `json:"user_id"`
	ShiftID        uint       `json:"shift_id"`
	Date           time.Time  `json:"date"`
	OfficeID       uint       `json:"office_id"`
	Status         string     `json:"status"`

	// Check-in
	TimeIn         *time.Time `json:"time_in"`
	LatIn          *float64   `json:"lat_in"`
	LngIn          *float64   `json:"lng_in"`
	ImageIn        string     `json:"image_in"`
	DistanceMeters *float64   `json:"distance_meters"`

	// Check-out
	TimeOut        *time.Time `json:"time_out"`
	LatOut         *float64   `json:"lat_out"`
	LngOut         *float64   `json:"lng_out"`
	ImageOut       string     `json:"image_out"`
	Duration       string     `json:"duration"`
}

func ToAttendanceDTO(a *models.Attendance) AttendanceDTO {
	return AttendanceDTO{
		ID:             a.ID,
		UserID:         a.UserID,
		ShiftID:        a.ShiftID,
		Date:           a.Date,
		OfficeID:       a.OfficeID,
		Status:         a.Status,
		TimeIn:         a.TimeIn,
		LatIn:          a.LatIn,
		LngIn:          a.LngIn,
		ImageIn:        a.ImageIn,
		DistanceMeters: a.DistanceMeters,
		TimeOut:        a.TimeOut,
		LatOut:         a.LatOut,
		LngOut:         a.LngOut,
		ImageOut:       a.ImageOut,
		Duration:       a.Duration,
	}
}

func ToAttendanceDTOList(attendances []models.Attendance) []AttendanceDTO {
	result := make([]AttendanceDTO, len(attendances))
	for i, a := range attendances {
		result[i] = ToAttendanceDTO(&a)
	}
	return result
}
