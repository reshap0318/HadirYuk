package dtos

import "github.com/reshap0318/hadirYuk/internal/models"

type LocationRequest struct {
	Name         string  `json:"name" validate:"required"`
	Address      string  `json:"address" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required"`
	Longitude    float64 `json:"longitude" validate:"required"`
	RadiusMeters int     `json:"radius_meters" validate:"required,min=50,max=500"`
	IsActive     bool    `json:"is_active"`
}

type LocationDTO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RadiusMeters int     `json:"radius_meters"`
	IsActive     bool    `json:"is_active"`
}

func ToLocationDTO(l *models.OfficeLocation) LocationDTO {
	return LocationDTO{
		ID:           l.ID,
		Name:         l.Name,
		Address:      l.Address,
		Latitude:     l.Latitude,
		Longitude:    l.Longitude,
		RadiusMeters: l.RadiusMeters,
		IsActive:     l.IsActive,
	}
}

func ToLocationDTOList(locations []models.OfficeLocation) []LocationDTO {
	result := make([]LocationDTO, len(locations))
	for i, l := range locations {
		result[i] = ToLocationDTO(&l)
	}
	return result
}
