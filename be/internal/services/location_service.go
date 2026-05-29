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

func (s *Services) LocationCreate(ctx context.Context, req dtos.LocationRequest) (*dtos.LocationDTO, error) {
	s.Logger.LogStart("LocationCreate", "Creating location: %s", req.Name)

	location := &models.OfficeLocation{
		Name:         req.Name,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: req.RadiusMeters,
		IsActive:     req.IsActive,
	}

	var result *models.OfficeLocation
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.OfficeLocation.Create(tx, location)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "success",
			Title:   "Location Created",
			Message: fmt.Sprintf("New location created: %s", result.Name),
			Data: map[string]interface{}{
				"id":   result.ID,
				"name": result.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("LocationCreate", "Failed to create location: %v", err)
		return nil, err
	}

	result = res.(*models.OfficeLocation)
	dto := dtos.ToLocationDTO(result)
	s.Logger.LogEnd("LocationCreate", "Location created: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

func (s *Services) LocationGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.LocationDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "ASC"
	}

	result, err := s.repo.OfficeLocation.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	locationDTOs := make([]dtos.LocationDTO, len(result.Data))
	for i, l := range result.Data {
		locationDTOs[i] = dtos.ToLocationDTO(&l)
	}

	return &repositories.PagedResult[dtos.LocationDTO]{
		Data:       locationDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *Services) LocationGetAllUnpaginated(ctx context.Context) ([]dtos.LocationDTO, error) {
	locations, err := s.repo.OfficeLocation.FindAll(nil)
	if err != nil {
		return nil, err
	}

	return dtos.ToLocationDTOList(locations), nil
}

func (s *Services) LocationGetByID(ctx context.Context, id uint) (*dtos.LocationDTO, error) {
	location, err := s.repo.OfficeLocation.FindByID(nil, id)
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToLocationDTO(location)
	return &dto, nil
}

func (s *Services) LocationUpdate(ctx context.Context, id uint, req dtos.LocationRequest) (*dtos.LocationDTO, error) {
	s.Logger.LogStart("LocationUpdate", "Updating location ID: %d", id)

	location := &models.OfficeLocation{ID: id}
	location.Name = req.Name
	location.Address = req.Address
	location.Latitude = req.Latitude
	location.Longitude = req.Longitude
	location.RadiusMeters = req.RadiusMeters
	location.IsActive = req.IsActive

	var result *models.OfficeLocation
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.OfficeLocation.Update(tx, &models.OfficeLocation{ID: id}, location)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "info",
			Title:   "Location Updated",
			Message: fmt.Sprintf("Location updated: %s", result.Name),
			Data: map[string]interface{}{
				"id":   result.ID,
				"name": result.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("LocationUpdate", "Failed to update location: %v", err)
		return nil, err
	}

	result = res.(*models.OfficeLocation)
	dto := dtos.ToLocationDTO(result)
	s.Logger.LogEnd("LocationUpdate", "Location updated: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

func (s *Services) LocationDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("LocationDelete", "Deleting location ID: %d", id)

	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		location, err := s.repo.OfficeLocation.FindByID(nil, id)
		if err != nil {
			return err
		}

		_, err = s.repo.OfficeLocation.Delete(tx, id)
		if err != nil {
			return err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "warning",
			Title:   "Location Deleted",
			Message: fmt.Sprintf("Location deleted: %s", location.Name),
			Data: map[string]interface{}{
				"id":   location.ID,
				"name": location.Name,
			},
		})

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("LocationDelete", "Failed to delete location: %v", err)
		return err
	}

	s.Logger.LogEnd("LocationDelete", "Location deleted: ID: %d", id)
	return nil
}
