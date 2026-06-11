package dtos

import (
	"time"

	"github.com/reshap0318/hadirYuk/internal/models"
)

type QRCodeGenerateRequest struct {
	OfficeID uint   `json:"office" validate:"required"`
	EndDate  string `json:"end_date" validate:"required"`
	EndTime  string `json:"end_time" validate:"required"`
}

type QRCodeDTO struct {
	ID        uint      `json:"id"`
	OfficeID  uint      `json:"office_id"`
	Office    *LocationDTO `json:"office,omitempty"`
	CodeValue string    `json:"code_value"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
	CreatedBy uint      `json:"created_by"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type QRCodeGenerateResponse struct {
	ID        uint      `json:"id"`
	OfficeID  uint      `json:"office_id"`
	Office    *LocationDTO `json:"office,omitempty"`
	CodeValue string    `json:"code_value"`
	Signature string    `json:"signature"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func ToQRCodeDTO(q *models.QRCode) QRCodeDTO {
	dto := QRCodeDTO{
		ID:        q.ID,
		OfficeID:  q.OfficeID,
		CodeValue: q.CodeValue,
		ExpiresAt: q.ExpiresAt,
		IsActive:  q.IsActive,
		CreatedBy: q.CreatedBy,
		RevokedAt: q.RevokedAt,
		CreatedAt: q.CreatedAt,
	}
	if q.Office.ID != 0 {
		officeDTO := ToLocationDTO(&q.Office)
		dto.Office = &officeDTO
	}
	return dto
}

func ToQRCodeDTOList(qrCodes []models.QRCode) []QRCodeDTO {
	result := make([]QRCodeDTO, len(qrCodes))
	for i, q := range qrCodes {
		result[i] = ToQRCodeDTO(&q)
	}
	return result
}

func ToQRCodeGenerateResponse(q *models.QRCode) QRCodeGenerateResponse {
	resp := QRCodeGenerateResponse{
		ID:        q.ID,
		OfficeID:  q.OfficeID,
		CodeValue: q.CodeValue,
		Signature: q.Signature,
		ExpiresAt: q.ExpiresAt,
		IsActive:  q.IsActive,
		CreatedAt: q.CreatedAt,
	}
	if q.Office.ID != 0 {
		officeDTO := ToLocationDTO(&q.Office)
		resp.Office = &officeDTO
	}
	return resp
}
