package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
)

func (s *Services) QRCodeGenerate(ctx context.Context, req *dtos.QRCodeGenerateRequest) (*dtos.QRCodeGenerateResponse, error) {
	s.Logger.LogStart("QRCodeGenerate", "Generating QR code for office %d", req.OfficeID)

	expiresAt, err := time.Parse("2006-01-02 15:04", req.EndDate+" "+req.EndTime)
	if err != nil {
		s.Logger.LogEndWithError("QRCodeGenerate", "Invalid date/time format: %v", err)
		return nil, fmt.Errorf("format tanggal atau waktu tidak valid, gunakan YYYY-MM-DD dan HH:MM")
	}

	var result *models.QRCode
	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		office, err := s.repo.OfficeLocation.FindByID(nil, req.OfficeID)
		if err != nil {
			return helpers.ErrNotFound
		}

		codeValue := uuid.New().String()
		signature := s.generateSignature(codeValue, expiresAt)

		callerID := helpers.GetCallerID(ctx)

		qrCode := &models.QRCode{
			OfficeID:  office.ID,
			CodeValue: codeValue,
			Signature: signature,
			ExpiresAt: expiresAt,
			IsActive:  true,
			CreatedBy: callerID,
		}

		var err2 error
		result, err2 = s.repo.QRCode.Create(tx, qrCode)
		if err2 != nil {
			return err2
		}

		result.Office = *office

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "success",
			Title:   "QR Code Generated",
			Message: fmt.Sprintf("New QR code generated for office: %s", office.Name),
			Data: map[string]interface{}{
				"id":         result.ID,
				"office_id":  result.OfficeID,
				"expires_at": result.ExpiresAt,
			},
		})

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("QRCodeGenerate", "Failed: %v", err)
		return nil, err
	}

	s.Logger.LogEnd("QRCodeGenerate", "QR code generated: %s (ID: %d)", result.CodeValue, result.ID)
	resp := dtos.ToQRCodeGenerateResponse(result)
	return &resp, nil
}

func (s *Services) QRCodeGetAll(ctx context.Context) ([]dtos.QRCodeDTO, error) {
	qrCodes, err := s.repo.QRCode.FindByField(nil, &models.QRCode{IsActive: true}, "Office")
	if err != nil {
		return nil, err
	}

	activeQRCodes := make([]models.QRCode, 0, len(qrCodes))
	now := time.Now()
	for _, q := range qrCodes {
		if q.ExpiresAt.After(now) && q.RevokedAt == nil {
			activeQRCodes = append(activeQRCodes, q)
		}
	}

	result := make([]dtos.QRCodeDTO, len(activeQRCodes))
	for i, q := range activeQRCodes {
		result[i] = dtos.ToQRCodeDTO(&q)
	}
	return result, nil
}

func (s *Services) QRCodeRevoke(ctx context.Context, id uint) (*dtos.QRCodeDTO, error) {
	s.Logger.LogStart("QRCodeRevoke", "Revoking QR code %d", id)

	var result *models.QRCode
	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		qrCode, err := s.repo.QRCode.FindByID(nil, id)
		if err != nil {
			return helpers.ErrNotFound
		}

		if !qrCode.IsActive {
			return fmt.Errorf("qr code already revoked")
		}

		now := time.Now()
		result, err = s.repo.QRCode.Update(tx, &models.QRCode{ID: id}, &models.QRCode{
			IsActive:  false,
			RevokedAt: &now,
		})
		if err != nil {
			return err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "warning",
			Title:   "QR Code Revoked",
			Message: fmt.Sprintf("QR code %s has been revoked", qrCode.CodeValue),
			Data: map[string]interface{}{
				"id":         result.ID,
				"code_value": result.CodeValue,
			},
		})

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("QRCodeRevoke", "Failed: %v", err)
		return nil, err
	}

	s.Logger.LogEnd("QRCodeRevoke", "QR code %d revoked", id)
	dto := dtos.ToQRCodeDTO(result)
	return &dto, nil
}

func (s *Services) QRCodeValidate(ctx context.Context, codeValue string) (*models.QRCode, error) {
	s.Logger.LogStart("QRCodeValidate", "Validating QR code: %s", codeValue)

	qrCodes, err := s.repo.QRCode.FindByFieldMap(nil, map[string]interface{}{
		"code_value": codeValue,
		"is_active":  true,
	})
	if err != nil || len(qrCodes) == 0 {
		s.Logger.LogEndWithError("QRCodeValidate", "QR code not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	qrCode := &qrCodes[0]

	if qrCode.RevokedAt != nil {
		s.Logger.LogEndWithError("QRCodeValidate", "QR code revoked")
		return nil, fmt.Errorf("qr code sudah dicabut")
	}

	if time.Now().After(qrCode.ExpiresAt) {
		s.Logger.LogEndWithError("QRCodeValidate", "QR code expired")
		return nil, fmt.Errorf("qr code sudah expired")
	}

	expectedSignature := s.generateSignature(qrCode.CodeValue, qrCode.ExpiresAt)
	if !hmac.Equal([]byte(qrCode.Signature), []byte(expectedSignature)) {
		s.Logger.LogEndWithError("QRCodeValidate", "Invalid signature")
		return nil, fmt.Errorf("signature qr code tidak valid")
	}

	s.Logger.LogEnd("QRCodeValidate", "QR code valid: %s", codeValue)
	return qrCode, nil
}

func (s *Services) generateSignature(codeValue string, expiresAt time.Time) string {
	secret := helpers.GetEnv("QR_SECRET_KEY", "default-qr-secret-key")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(codeValue + expiresAt.Format(time.RFC3339)))
	return hex.EncodeToString(mac.Sum(nil))
}
