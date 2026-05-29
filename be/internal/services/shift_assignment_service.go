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

	// Parse dates first
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.Logger.LogEndWithError("ShiftAssignToUser", "Invalid start_date format: %v", err)
		return nil, &helpers.FieldError{Field: "start_date", Message: "Invalid date format, use YYYY-MM-DD"}
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsedEnd, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			s.Logger.LogEndWithError("ShiftAssignToUser", "Invalid end_date format: %v", err)
			return nil, &helpers.FieldError{Field: "end_date", Message: "Invalid date format, use YYYY-MM-DD"}
		}
		endDate = &parsedEnd
	}

	// Check for overlapping active assignment (same user + shift)
	overlapping, err := s.repo.UserShiftAssignment.FindOverlappingAssignments(nil, req.UserID, req.ShiftID, startDate, endDate, 0)
	if err != nil {
		s.Logger.LogEndWithError("ShiftAssignToUser", "Failed to check existing assignment: %v", err)
		return nil, err
	}
	if len(overlapping) > 0 {
		s.Logger.LogEndWithError("ShiftAssignToUser", "User already has an overlapping shift assignment for this period")
		return nil, &helpers.FieldError{Field: "start_date", Message: "User already has an overlapping shift assignment for this period"}
	}

	assignment := &models.UserShiftAssignment{
		UserID:    req.UserID,
		ShiftID:   req.ShiftID,
		StartDate: startDate,
		IsActive:  true,
		EndDate:   endDate,
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

	// Determine effective values for overlap check
	effectiveUserID := existing.UserID
	effectiveShiftID := existing.ShiftID
	effectiveStartDate := existing.StartDate
	effectiveEndDate := existing.EndDate

	if req.ShiftID != 0 {
		// Check new shift exists
		_, err := s.repo.Shift.FindByID(nil, req.ShiftID)
		if err != nil {
			s.Logger.LogEndWithError("ShiftUpdateAssignment", "New shift not found: %v", err)
			return nil, &helpers.FieldError{Field: "shift_id", Message: "Shift not found"}
		}
		updates["shift_id"] = req.ShiftID
		effectiveShiftID = req.ShiftID
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			s.Logger.LogEndWithError("ShiftUpdateAssignment", "Invalid start_date format: %v", err)
			return nil, &helpers.FieldError{Field: "start_date", Message: "Invalid date format, use YYYY-MM-DD"}
		}
		updates["start_date"] = startDate
		effectiveStartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			s.Logger.LogEndWithError("ShiftUpdateAssignment", "Invalid end_date format: %v", err)
			return nil, &helpers.FieldError{Field: "end_date", Message: "Invalid date format, use YYYY-MM-DD"}
		}
		updates["end_date"] = &endDate
		effectiveEndDate = &endDate
	} else if req.EndDate == "" && existing.EndDate != nil && req.StartDate != "" {
		// If end_date not provided in request but start_date changed, keep existing end_date
		// (already set above via effectiveEndDate)
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		dto := dtos.ToShiftAssignmentDTO(existing)
		s.Logger.LogEnd("ShiftUpdateAssignment", "No updates provided for assignment ID: %d", id)
		return &dto, nil
	}

	// Check for overlapping active assignment (same user + shift, exclude current record)
	// Only check if dates or shift changed
	if req.StartDate != "" || req.EndDate != "" || req.ShiftID != 0 {
		overlapping, err := s.repo.UserShiftAssignment.FindOverlappingAssignments(nil, effectiveUserID, effectiveShiftID, effectiveStartDate, effectiveEndDate, id)
		if err != nil {
			s.Logger.LogEndWithError("ShiftUpdateAssignment", "Failed to check existing assignment: %v", err)
			return nil, err
		}
		if len(overlapping) > 0 {
			s.Logger.LogEndWithError("ShiftUpdateAssignment", "User already has an overlapping shift assignment for this period")
			return nil, &helpers.FieldError{Field: "start_date", Message: "User already has an overlapping shift assignment for this period"}
		}
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
