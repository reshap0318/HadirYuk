package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

type NotificationCreateParams struct {
	Type    string
	Title   string
	Message string
	Data    map[string]interface{}
}

func (s *Services) NotificationCreate(ctx context.Context, params *NotificationCreateParams) error {
	userID := helpers.GetCallerID(ctx)
	s.Logger.LogStart("NotificationCreate", "Creating notification for user %d: %s", userID, params.Title)

	var dataJSON string
	if params.Data != nil {
		dataBytes, err := json.Marshal(params.Data)
		if err != nil {
			s.Logger.LogEndWithError("NotificationCreate", "Failed to marshal data: %v", err)
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
		dataJSON = string(dataBytes)
	}

	notification := &models.Notification{
		UserID:  userID,
		Type:    params.Type,
		Title:   params.Title,
		Message: params.Message,
		Data:    dataJSON,
	}

	_, err := s.repo.Notification.Create(nil, notification)
	if err != nil {
		s.Logger.LogEndWithError("NotificationCreate", "Failed to create notification: %v", err)
		return fmt.Errorf("failed to create notification: %w", err)
	}

	s.Logger.LogEnd("NotificationCreate", "Notification created for user %d", userID)
	return nil
}

func (s *Services) NotificationGetAllWithFilters(ctx context.Context, opts *repositories.QueryOptions, isRead *bool, notifType string) (*repositories.PagedResult[dtos.NotificationDTO], error) {
	s.Logger.LogStart("NotificationGetAllWithFilters", "Fetching notifications with filters")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		return nil, helpers.ErrInvalidToken
	}

	result, err := s.repo.Notification.FindByUserIDAndFilters(userID, opts, isRead, notifType)
	if err != nil {
		s.Logger.LogEndWithError("NotificationGetAllWithFilters", "Failed to fetch notifications: %v", err)
		return nil, err
	}

	dtoList := dtos.ToNotificationDTOList(result.Data)

	s.Logger.LogEnd("NotificationGetAllWithFilters", "Fetched %d notifications for user %d", len(dtoList), userID)
	return &repositories.PagedResult[dtos.NotificationDTO]{
		Data:       dtoList,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *Services) NotificationGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.NotificationDTO], error) {
	s.Logger.LogStart("NotificationGetAll", "Fetching notifications")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		return nil, helpers.ErrInvalidToken
	}

	result, err := s.repo.Notification.FindByUserIDWithOpts(userID, opts)
	if err != nil {
		s.Logger.LogEndWithError("NotificationGetAll", "Failed to fetch notifications: %v", err)
		return nil, err
	}

	dtoList := dtos.ToNotificationDTOList(result.Data)

	s.Logger.LogEnd("NotificationGetAll", "Fetched %d notifications for user %d", len(dtoList), userID)
	return &repositories.PagedResult[dtos.NotificationDTO]{
		Data:       dtoList,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *Services) NotificationGetByID(ctx context.Context, id uint) (*dtos.NotificationDTO, error) {
	s.Logger.LogStart("NotificationGetByID", "Fetching notification %d", id)

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		return nil, helpers.ErrInvalidToken
	}

	notification, err := s.repo.Notification.FindByID(nil, id)
	if err != nil {
		s.Logger.LogEndWithError("NotificationGetByID", "Failed to fetch notification: %v", err)
		return nil, err
	}

	if notification.UserID != userID {
		return nil, helpers.ErrForbidden
	}

	dto := dtos.ToNotificationDTO(notification)

	s.Logger.LogEnd("NotificationGetByID", "Fetched notification %d", id)
	return &dto, nil
}

func (s *Services) NotificationMarkAsRead(ctx context.Context, id uint) error {
	s.Logger.LogStart("NotificationMarkAsRead", "Marking notification %d as read", id)

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		return helpers.ErrInvalidToken
	}

	err := s.repo.Notification.MarkAsRead(id, userID)
	if err != nil {
		s.Logger.LogEndWithError("NotificationMarkAsRead", "Failed to mark notification as read: %v", err)
		return err
	}

	s.Logger.LogEnd("NotificationMarkAsRead", "Notification %d marked as read", id)
	return nil
}

func (s *Services) NotificationMarkAllAsRead(ctx context.Context) error {
	s.Logger.LogStart("NotificationMarkAllAsRead", "Marking all notifications as read")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		return helpers.ErrInvalidToken
	}

	err := s.repo.Notification.MarkAllAsRead(userID)
	if err != nil {
		s.Logger.LogEndWithError("NotificationMarkAllAsRead", "Failed to mark all notifications as read: %v", err)
		return err
	}

	s.Logger.LogEnd("NotificationMarkAllAsRead", "All notifications marked as read for user %d", userID)
	return nil
}

func (s *Services) NotificationCountUnread(ctx context.Context) (int64, error) {
	s.Logger.LogStart("NotificationCountUnread", "Counting unread notifications")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		return 0, helpers.ErrInvalidToken
	}

	count, err := s.repo.Notification.CountUnread(userID)
	if err != nil {
		s.Logger.LogEndWithError("NotificationCountUnread", "Failed to count unread notifications: %v", err)
		return 0, err
	}

	s.Logger.LogEnd("NotificationCountUnread", "User %d has %d unread notifications", userID, count)
	return count, nil
}

func (s *Services) NotificationDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("NotificationDelete", "Deleting notification %d", id)

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		return helpers.ErrInvalidToken
	}

	notification, err := s.repo.Notification.FindByID(nil, id)
	if err != nil {
		s.Logger.LogEndWithError("NotificationDelete", "Failed to fetch notification: %v", err)
		return err
	}

	if notification.UserID != userID {
		return helpers.ErrForbidden
	}

	_, err = s.repo.Notification.Delete(nil, id)
	if err != nil {
		s.Logger.LogEndWithError("NotificationDelete", "Failed to delete notification: %v", err)
		return err
	}

	s.Logger.LogEnd("NotificationDelete", "Notification %d deleted", id)
	return nil
}
