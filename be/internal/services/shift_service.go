package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

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
		TotalHours:    req.TotalHours,
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

	exists, err := s.repo.Shift.Exists(nil, map[string]interface{}{"name": req.Name})
	if err != nil {
		s.Logger.LogEndWithError("ShiftUpdate", "Failed to check duplicate: %v", err)
		return nil, err
	}
	if exists {
		shift, _ := s.repo.Shift.FindByID(nil, id)
		if shift == nil || shift.Name != req.Name {
			return nil, &helpers.FieldError{Field: "name", Message: "Shift name already exists"}
		}
	}

	shift := &models.Shift{ID: id}
	if req.Name != "" {
		shift.Name = req.Name
	}
	if req.StartTime != "" {
		shift.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		shift.EndTime = req.EndTime
	}
	shift.BreakDuration = req.BreakDuration
	if req.ColorCode != "" {
		shift.ColorCode = req.ColorCode
	}
	shift.TotalHours = req.TotalHours

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

func (s *Services) ShiftAssignToUser(ctx context.Context, req dtos.ShiftAssignmentRequest) (*dtos.ShiftAssignmentDTO, error) {
	s.Logger.LogStart("ShiftAssignToUser", "Assigning shift to user ID: %d, shift ID: %d", req.UserID, req.ShiftID)

	// Check user exists
	user, err := s.repo.User.FindByID(nil, req.UserID)
	if err != nil {
		s.Logger.LogEndWithError("ShiftAssignToUser", "User not found: %v", err)
		return nil, &helpers.FieldError{Field: "user_id", Message: "User not found"}
	}

	// Check shift exists
	shift, err := s.repo.Shift.FindByID(nil, req.ShiftID)
	if err != nil {
		s.Logger.LogEndWithError("ShiftAssignToUser", "Shift not found: %v", err)
		return nil, &helpers.FieldError{Field: "shift_id", Message: "Shift not found"}
	}

	// Check no overlapping active assignment
	exists, err := s.repo.UserShiftAssignment.Exists(nil, map[string]interface{}{
		"user_id":   req.UserID,
		"is_active": true,
	})
	if err != nil {
		s.Logger.LogEndWithError("ShiftAssignToUser", "Failed to check existing assignment: %v", err)
		return nil, err
	}
	if exists {
		s.Logger.LogEndWithError("ShiftAssignToUser", "User already has an active shift assignment")
		return nil, &helpers.FieldError{Field: "user_id", Message: "User already has an active shift assignment"}
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.Logger.LogEndWithError("ShiftAssignToUser", "Invalid start_date format: %v", err)
		return nil, &helpers.FieldError{Field: "start_date", Message: "Invalid date format, use YYYY-MM-DD"}
	}

	assignment := &models.UserShiftAssignment{
		UserID:    req.UserID,
		ShiftID:   req.ShiftID,
		StartDate: startDate,
		IsActive:  true,
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			s.Logger.LogEndWithError("ShiftAssignToUser", "Invalid end_date format: %v", err)
			return nil, &helpers.FieldError{Field: "end_date", Message: "Invalid date format, use YYYY-MM-DD"}
		}
		assignment.EndDate = &endDate
	}

	var result *models.UserShiftAssignment
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.UserShiftAssignment.Create(tx, assignment)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "success",
			Title:   "Shift Assigned",
			Message: fmt.Sprintf("Shift %s assigned to %s", shift.Name, user.Name),
			Data: map[string]interface{}{
				"user_id":  req.UserID,
				"shift_id": req.ShiftID,
				"user":     user.Name,
				"shift":    shift.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ShiftAssignToUser", "Failed to assign shift: %v", err)
		return nil, err
	}

	result = res.(*models.UserShiftAssignment)
	dto := dtos.ToShiftAssignmentDTO(result)
	s.Logger.LogEnd("ShiftAssignToUser", "Shift assigned: user ID %d, shift ID %d", req.UserID, req.ShiftID)
	return &dto, nil
}

func (s *Services) ShiftGetUserAssignments(ctx context.Context, userID uint, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.ShiftAssignmentDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "DESC"
	}

	result, err := s.repo.UserShiftAssignment.FindByUserIDWithHistory(nil, userID, opts, "User", "Shift")
	if err != nil {
		return nil, err
	}

	assignmentDTOs := make([]dtos.ShiftAssignmentDTO, len(result.Data))
	for i, a := range result.Data {
		assignmentDTOs[i] = dtos.ToShiftAssignmentDTO(&a)
	}

	return &repositories.PagedResult[dtos.ShiftAssignmentDTO]{
		Data:       assignmentDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *Services) ShiftGetAllAssignments(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.ShiftAssignmentDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "DESC"
	}
	if len(opts.Preloads) == 0 {
		opts.Preloads = []string{"User", "Shift"}
	}

	result, err := s.repo.UserShiftAssignment.FindAllWithSearch(nil, opts)
	if err != nil {
		return nil, err
	}

	assignmentDTOs := make([]dtos.ShiftAssignmentDTO, len(result.Data))
	for i, a := range result.Data {
		assignmentDTOs[i] = dtos.ToShiftAssignmentDTO(&a)
	}

	return &repositories.PagedResult[dtos.ShiftAssignmentDTO]{
		Data:       assignmentDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *Services) ShiftGetUserActiveAssignment(ctx context.Context, userID uint) (*dtos.ShiftAssignmentDTO, error) {
	assignment, err := s.repo.UserShiftAssignment.FindByUserID(nil, userID, "User", "Shift")
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToShiftAssignmentDTO(assignment)
	return &dto, nil
}

func (s *Services) ShiftUpdateAssignment(ctx context.Context, id uint, req dtos.ShiftAssignmentUpdateRequest) (*dtos.ShiftAssignmentDTO, error) {
	s.Logger.LogStart("ShiftUpdateAssignment", "Updating assignment ID: %d", id)

	// Check assignment exists
	existing, err := s.repo.UserShiftAssignment.FindByID(nil, id, "User", "Shift")
	if err != nil {
		s.Logger.LogEndWithError("ShiftUpdateAssignment", "Assignment not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	updates := map[string]interface{}{}

	if req.ShiftID != 0 {
		// Check new shift exists
		_, err := s.repo.Shift.FindByID(nil, req.ShiftID)
		if err != nil {
			s.Logger.LogEndWithError("ShiftUpdateAssignment", "New shift not found: %v", err)
			return nil, &helpers.FieldError{Field: "shift_id", Message: "Shift not found"}
		}
		updates["shift_id"] = req.ShiftID
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			s.Logger.LogEndWithError("ShiftUpdateAssignment", "Invalid start_date format: %v", err)
			return nil, &helpers.FieldError{Field: "start_date", Message: "Invalid date format, use YYYY-MM-DD"}
		}
		updates["start_date"] = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			s.Logger.LogEndWithError("ShiftUpdateAssignment", "Invalid end_date format: %v", err)
			return nil, &helpers.FieldError{Field: "end_date", Message: "Invalid date format, use YYYY-MM-DD"}
		}
		updates["end_date"] = &endDate
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		dto := dtos.ToShiftAssignmentDTO(existing)
		s.Logger.LogEnd("ShiftUpdateAssignment", "No updates provided for assignment ID: %d", id)
		return &dto, nil
	}

	var result *models.UserShiftAssignment
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.UserShiftAssignment.UpdateMap(tx, &models.UserShiftAssignment{ID: id}, updates)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "info",
			Title:   "Shift Assignment Updated",
			Message: fmt.Sprintf("Shift assignment updated for %s", existing.User.Name),
			Data: map[string]interface{}{
				"id":        id,
				"user_id":   existing.UserID,
				"user":      existing.User.Name,
				"is_active": result.IsActive,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ShiftUpdateAssignment", "Failed to update assignment: %v", err)
		return nil, err
	}

	result = res.(*models.UserShiftAssignment)
	dto := dtos.ToShiftAssignmentDTO(result)
	s.Logger.LogEnd("ShiftUpdateAssignment", "Assignment updated: ID %d", id)
	return &dto, nil
}

func (s *Services) ShiftDeleteAssignment(ctx context.Context, id uint) error {
	s.Logger.LogStart("ShiftDeleteAssignment", "Deleting assignment ID: %d", id)

	existing, err := s.repo.UserShiftAssignment.FindByID(nil, id, "User", "Shift")
	if err != nil {
		s.Logger.LogEndWithError("ShiftDeleteAssignment", "Assignment not found: %v", err)
		return helpers.ErrNotFound
	}

	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.UserShiftAssignment.Delete(tx, id)
		if err != nil {
			return err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "warning",
			Title:   "Shift Assignment Deleted",
			Message: fmt.Sprintf("Shift assignment deleted for %s", existing.User.Name),
			Data: map[string]interface{}{
				"id":      id,
				"user_id": existing.UserID,
				"user":    existing.User.Name,
			},
		})

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ShiftDeleteAssignment", "Failed to delete assignment: %v", err)
		return err
	}

	s.Logger.LogEnd("ShiftDeleteAssignment", "Assignment deleted: ID %d", id)
	return nil
}
