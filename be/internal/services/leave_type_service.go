package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

func (s *Services) LeaveTypeCreate(ctx context.Context, req dtos.LeaveTypeRequest) (*dtos.LeaveTypeDTO, error) {
	s.Logger.LogStart("LeaveTypeCreate", "Creating leave type: %s", req.Name)

	exists, err := s.repo.LeaveType.Exists(nil, map[string]interface{}{"name": req.Name})
	if err != nil {
		s.Logger.LogEndWithError("LeaveTypeCreate", "Failed to check duplicate: %v", err)
		return nil, err
	}
	if exists {
		return nil, &helpers.FieldError{Field: "name", Message: "Leave type name already exists"}
	}

	leaveType := &models.LeaveType{
		Name:        req.Name,
		Description: req.Description,
		DefaultDays: req.DefaultDays,
		IsPaid:      req.IsPaid,
	}

	var result *models.LeaveType
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.LeaveType.Create(tx, leaveType)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "success",
			Title:   "Leave Type Created",
			Message: fmt.Sprintf("New leave type created: %s", result.Name),
			Data: map[string]interface{}{
				"id":   result.ID,
				"name": result.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("LeaveTypeCreate", "Failed to create leave type: %v", err)
		return nil, err
	}

	result = res.(*models.LeaveType)
	dto := dtos.ToLeaveTypeDTO(result)
	s.Logger.LogEnd("LeaveTypeCreate", "Leave type created: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

func (s *Services) LeaveTypeGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.LeaveTypeDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "ASC"
	}

	result, err := s.repo.LeaveType.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	leaveTypeDTOs := make([]dtos.LeaveTypeDTO, len(result.Data))
	for i, l := range result.Data {
		leaveTypeDTOs[i] = dtos.ToLeaveTypeDTO(&l)
	}

	return &repositories.PagedResult[dtos.LeaveTypeDTO]{
		Data:       leaveTypeDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *Services) LeaveTypeGetAllUnpaginated(ctx context.Context) ([]dtos.LeaveTypeDTO, error) {
	leaveTypes, err := s.repo.LeaveType.FindAll(nil)
	if err != nil {
		return nil, err
	}

	return dtos.ToLeaveTypeDTOList(leaveTypes), nil
}

func (s *Services) LeaveTypeGetByID(ctx context.Context, id uint) (*dtos.LeaveTypeDTO, error) {
	leaveType, err := s.repo.LeaveType.FindByID(nil, id)
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToLeaveTypeDTO(leaveType)
	return &dto, nil
}

func (s *Services) LeaveTypeUpdate(ctx context.Context, id uint, req dtos.LeaveTypeRequest) (*dtos.LeaveTypeDTO, error) {
	s.Logger.LogStart("LeaveTypeUpdate", "Updating leave type ID: %d", id)

	exists, err := s.repo.LeaveType.Exists(nil, map[string]interface{}{"name": req.Name})
	if err != nil {
		s.Logger.LogEndWithError("LeaveTypeUpdate", "Failed to check duplicate: %v", err)
		return nil, err
	}
	if exists {
		leaveType, _ := s.repo.LeaveType.FindByID(nil, id)
		if leaveType == nil || leaveType.Name != req.Name {
			return nil, &helpers.FieldError{Field: "name", Message: "Leave type name already exists"}
		}
	}

	leaveType := &models.LeaveType{ID: id}
	leaveType.Name = req.Name
	leaveType.Description = req.Description
	leaveType.DefaultDays = req.DefaultDays
	leaveType.IsPaid = req.IsPaid

	var result *models.LeaveType
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.LeaveType.Update(tx, &models.LeaveType{ID: id}, leaveType)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "info",
			Title:   "Leave Type Updated",
			Message: fmt.Sprintf("Leave type updated: %s", result.Name),
			Data: map[string]interface{}{
				"id":   result.ID,
				"name": result.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("LeaveTypeUpdate", "Failed to update leave type: %v", err)
		return nil, err
	}

	result = res.(*models.LeaveType)
	dto := dtos.ToLeaveTypeDTO(result)
	s.Logger.LogEnd("LeaveTypeUpdate", "Leave type updated: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

func (s *Services) LeaveTypeDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("LeaveTypeDelete", "Deleting leave type ID: %d", id)

	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		leaveType, err := s.repo.LeaveType.FindByID(nil, id)
		if err != nil {
			return err
		}

		_, err = s.repo.LeaveType.Delete(tx, id)
		if err != nil {
			return err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "warning",
			Title:   "Leave Type Deleted",
			Message: fmt.Sprintf("Leave type deleted: %s", leaveType.Name),
			Data: map[string]interface{}{
				"id":   leaveType.ID,
				"name": leaveType.Name,
			},
		})

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("LeaveTypeDelete", "Failed to delete leave type: %v", err)
		return err
	}

	s.Logger.LogEnd("LeaveTypeDelete", "Leave type deleted: ID: %d", id)
	return nil
}
