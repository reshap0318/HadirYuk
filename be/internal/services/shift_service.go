package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

// calculateTotalHours computes total working hours from start time, end time, and break duration.
// Time format is "HH:MM". Returns hours as float64 (e.g., 7.5 for 7 hours 30 minutes).
func calculateTotalHours(startTime, endTime string, breakDuration int) float64 {
	startParts := strings.Split(startTime, ":")
	endParts := strings.Split(endTime, ":")

	startHour, _ := strconv.Atoi(startParts[0])
	startMin, _ := strconv.Atoi(startParts[1])
	endHour, _ := strconv.Atoi(endParts[0])
	endMin, _ := strconv.Atoi(endParts[1])

	startTotalMin := startHour*60 + startMin
	endTotalMin := endHour*60 + endMin

	// Handle overnight shifts (e.g., 22:00 to 06:00)
	durationMin := endTotalMin - startTotalMin
	if durationMin < 0 {
		durationMin += 24 * 60
	}

	workMin := durationMin - breakDuration
	if workMin < 0 {
		workMin = 0
	}

	return float64(workMin) / 60.0
}

func (s *Services) ShiftCreate(ctx context.Context, req dtos.ShiftRequest) (*dtos.ShiftDTO, error) {
	s.Logger.LogStart("ShiftCreate", "Creating shift: %s", req.Name)

	exists, err := s.repo.Shift.Exists(nil, map[string]interface{}{"name": req.Name})
	if err != nil {
		s.Logger.LogEndWithError("ShiftCreate", "Failed to check duplicate: %v", err)
		return nil, err
	}
	if exists {
		return nil, &helpers.FieldError{Field: "name", Message: "Shift name already exists"}
	}

	shift := &models.Shift{
		Name:          req.Name,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		BreakDuration: req.BreakDuration,
		ColorCode:     req.ColorCode,
		TotalHours:    calculateTotalHours(req.StartTime, req.EndTime, req.BreakDuration),
	}

	var result *models.Shift
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.Shift.Create(tx, shift)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "success",
			Title:   "Shift Created",
			Message: fmt.Sprintf("New shift created: %s", result.Name),
			Data: map[string]interface{}{
				"id":   result.ID,
				"name": result.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ShiftCreate", "Failed to create shift: %v", err)
		return nil, err
	}

	result = res.(*models.Shift)
	dto := dtos.ToShiftDTO(result)
	s.Logger.LogEnd("ShiftCreate", "Shift created: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

func (s *Services) ShiftGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.ShiftDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "ASC"
	}

	result, err := s.repo.Shift.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	shiftDTOs := make([]dtos.ShiftDTO, len(result.Data))
	for i, sh := range result.Data {
		shiftDTOs[i] = dtos.ToShiftDTO(&sh)
	}

	return &repositories.PagedResult[dtos.ShiftDTO]{
		Data:       shiftDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *Services) ShiftGetAllUnpaginated(ctx context.Context) ([]dtos.ShiftDTO, error) {
	shifts, err := s.repo.Shift.FindAll(nil)
	if err != nil {
		return nil, err
	}

	return dtos.ToShiftDTOList(shifts), nil
}

func (s *Services) ShiftGetByID(ctx context.Context, id uint) (*dtos.ShiftDTO, error) {
	shift, err := s.repo.Shift.FindByID(nil, id)
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToShiftDTO(shift)
	return &dto, nil
}

func (s *Services) ShiftUpdate(ctx context.Context, id uint, req dtos.ShiftRequest) (*dtos.ShiftDTO, error) {
	s.Logger.LogStart("ShiftUpdate", "Updating shift ID: %d", id)

	_, err := s.repo.Shift.FindByID(nil, id)
	if err != nil {
		s.Logger.LogEndWithError("ShiftUpdate", "Shift not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	exists, err := s.repo.Shift.Exists(nil, map[string]interface{}{"name": req.Name, "id <>": id})
	if err != nil {
		s.Logger.LogEndWithError("ShiftUpdate", "Failed to check duplicate: %v", err)
		return nil, err
	}
	if exists {
		return nil, &helpers.FieldError{Field: "name", Message: "Shift name already exists"}
	}

	shift := &models.Shift{
		ID:            id,
		Name:          req.Name,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		BreakDuration: req.BreakDuration,
		ColorCode:     req.ColorCode,
		TotalHours:    calculateTotalHours(req.StartTime, req.EndTime, req.BreakDuration),
	}

	var result *models.Shift
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.Shift.Update(tx, &models.Shift{ID: id}, shift)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "info",
			Title:   "Shift Updated",
			Message: fmt.Sprintf("Shift updated: %s", result.Name),
			Data: map[string]interface{}{
				"id":   result.ID,
				"name": result.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ShiftUpdate", "Failed to update shift: %v", err)
		return nil, err
	}

	result = res.(*models.Shift)
	dto := dtos.ToShiftDTO(result)
	s.Logger.LogEnd("ShiftUpdate", "Shift updated: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

func (s *Services) ShiftDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("ShiftDelete", "Deleting shift ID: %d", id)

	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		shift, err := s.repo.Shift.FindByID(nil, id)
		if err != nil {
			return err
		}

		_, err = s.repo.Shift.Delete(tx, id)
		if err != nil {
			return err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "warning",
			Title:   "Shift Deleted",
			Message: fmt.Sprintf("Shift deleted: %s", shift.Name),
			Data: map[string]interface{}{
				"id":   shift.ID,
				"name": shift.Name,
			},
		})

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ShiftDelete", "Failed to delete shift: %v", err)
		return err
	}

	s.Logger.LogEnd("ShiftDelete", "Shift deleted: ID: %d", id)
	return nil
}
