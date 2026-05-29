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

// ShiftAssignToUser assigns a shift to a user.
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

// ShiftGetUserAssignments returns paginated shift assignments for a specific user.
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

// ShiftGetAllAssignments returns paginated shift assignments with search.
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

// ShiftGetUserActiveAssignment returns the active shift assignment for a user.
func (s *Services) ShiftGetUserActiveAssignment(ctx context.Context, userID uint) (*dtos.ShiftAssignmentDTO, error) {
	assignment, err := s.repo.UserShiftAssignment.FindByUserID(nil, userID, "User", "Shift")
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToShiftAssignmentDTO(assignment)
	return &dto, nil
}

// ShiftUpdateAssignment updates an existing shift assignment.
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

// ShiftDeleteAssignment deletes a shift assignment.
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
