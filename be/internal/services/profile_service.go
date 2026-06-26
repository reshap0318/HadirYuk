package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
)

// ProfileGet returns the authenticated user's profile.
func (s *Services) ProfileGet(ctx context.Context, userID uint) (*dtos.ProfileDTO, error) {
	s.Logger.LogStart("ProfileGet", "Fetching profile for user ID: %d", userID)

	user, err := s.repo.User.FindByID(nil, userID, "Roles", "Profile")
	if err != nil {
		s.Logger.LogEndWithError("ProfileGet", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToProfileDTO(user)
	s.Logger.LogEnd("ProfileGet", "Profile fetched for user: %s", dto.Email)
	return &dto, nil
}

// ProfileUpdate updates the authenticated user's profile.
func (s *Services) ProfileUpdate(ctx context.Context, userID uint, req dtos.ProfileUpdateRequest) (*dtos.ProfileDTO, error) {
	s.Logger.LogStart("ProfileUpdate", "Updating profile for user ID: %d", userID)

	existing, err := s.repo.User.FindByID(nil, userID, "Profile")
	if err != nil {
		s.Logger.LogEndWithError("ProfileUpdate", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	updates := map[string]interface{}{
		"name": req.Name,
	}

	if req.Avatar != "" {
		avatarPath, err := helpers.MoveFile(req.Avatar, "storage/tmp", "storage/avatars")
		if err != nil {
			s.Logger.LogStep("ProfileUpdate", "Failed to move avatar: %v", err)
		} else {
			if existing.Avatar != "" {
				helpers.DeleteFile(existing.Avatar)
			}
			updates["avatar"] = avatarPath
		}
	}

	profileUpdates := map[string]interface{}{}
	if req.Phone != "" {
		profileUpdates["phone"] = req.Phone
	}
	if req.Department != "" {
		profileUpdates["department"] = req.Department
	}
	if req.Position != "" {
		profileUpdates["position"] = req.Position
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		result, err := s.repo.User.UpdateMap(tx, &models.User{ID: userID}, updates)
		if err != nil {
			return nil, err
		}

		if len(profileUpdates) > 0 {
			if existing.Profile == nil {
				if _, err := s.repo.UserProfile.Create(tx, &models.UserProfile{UserID: userID}); err != nil {
					s.Logger.LogStep("ProfileUpdate", "Failed to create profile: %v", err)
					return nil, err
				}
			}
			if _, err := s.repo.UserProfile.UpdateMap(tx, &models.UserProfile{UserID: userID}, profileUpdates); err != nil {
				s.Logger.LogStep("ProfileUpdate", "Failed to update profile fields: %v", err)
				return nil, err
			}
		}

		reloaded, err := s.repo.User.FindByID(tx, result.ID, "Roles", "Profile")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ProfileUpdate", "Failed to update profile: %v", err)
		return nil, err
	}

	result := res.(*models.User)
	dto := dtos.ToProfileDTO(result)

	s.Access.Invalidate(userID)

	s.Logger.LogEnd("ProfileUpdate", "Profile updated for user: %s", dto.Email)
	return &dto, nil
}

// ProfileChangePassword changes the authenticated user's password.
func (s *Services) ProfileChangePassword(ctx context.Context, userID uint, req dtos.ChangePasswordRequest) error {
	s.Logger.LogStart("ProfileChangePassword", "Changing password for user ID: %d", userID)

	existing, err := s.repo.User.FindByID(nil, userID)
	if err != nil {
		s.Logger.LogEndWithError("ProfileChangePassword", "User not found: %v", err)
		return helpers.ErrNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		s.Logger.LogEndWithError("ProfileChangePassword", "Failed to hash password: %v", err)
		return err
	}

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.User.UpdateMap(tx, &models.User{ID: userID}, map[string]interface{}{
			"password": string(hashedPassword),
		})
		return err
	}); err != nil {
		s.Logger.LogEndWithError("ProfileChangePassword", "Failed to update password: %v", err)
		return err
	}

	// Invalidate Redis session cache — user must re-login
	if s.RedisClient.IsCacheAvailable() {
		sessionKey := fmt.Sprintf("session:%d", existing.ID)
		if err := s.RedisClient.Delete(sessionKey); err != nil {
			s.Logger.LogWarn("ProfileChangePassword", "Failed to delete session cache: %v", err)
		} else {
			s.Logger.LogStep("ProfileChangePassword", "Session cache invalidated")
		}
	}

	s.Access.Invalidate(userID)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "Password Changed",
		Message: "Your password has been changed successfully. Please re-login.",
		Data:    map[string]interface{}{"user_id": userID, "email": existing.Email},
	})

	s.Logger.LogEnd("ProfileChangePassword", "Password changed for user: %s", existing.Email)
	return nil
}
