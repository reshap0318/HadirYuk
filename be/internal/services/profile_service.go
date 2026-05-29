package services

import (
	"context"

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
			if result.Profile == nil {
				result.Profile = &models.UserProfile{UserID: userID}
				if _, err := s.repo.UserProfile.Create(tx, result.Profile); err != nil {
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

// ProfileUploadFacePhoto uploads a face photo for the authenticated user.
func (s *Services) ProfileUploadFacePhoto(ctx context.Context, userID uint, fileUUID string) (*dtos.ProfileDTO, error) {
	s.Logger.LogStart("ProfileUploadFacePhoto", "Uploading face photo for user ID: %d", userID)

	existing, err := s.repo.User.FindByID(nil, userID, "Profile")
	if err != nil {
		s.Logger.LogEndWithError("ProfileUploadFacePhoto", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	// Move file from tmp to face-photos directory
	facePhotoPath, err := helpers.MoveFile(fileUUID, "storage/tmp", "storage/face-photos")
	if err != nil {
		s.Logger.LogEndWithError("ProfileUploadFacePhoto", "Failed to move face photo: %v", err)
		return nil, err
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		// Delete old face photo if exists
		if existing.Profile != nil && existing.Profile.FacePhotoURL != "" {
			helpers.DeleteFile(existing.Profile.FacePhotoURL)
		}

		// Update or create profile with face photo
		profileUpdates := map[string]interface{}{
			"face_photo_url": facePhotoPath,
		}

		if existing.Profile == nil {
			profile := &models.UserProfile{
				UserID:       userID,
				FacePhotoURL: facePhotoPath,
			}
			if _, err := s.repo.UserProfile.Create(tx, profile); err != nil {
				s.Logger.LogStep("ProfileUploadFacePhoto", "Failed to create profile: %v", err)
				return nil, err
			}
		} else {
			if _, err := s.repo.UserProfile.UpdateMap(tx, &models.UserProfile{UserID: userID}, profileUpdates); err != nil {
				s.Logger.LogStep("ProfileUploadFacePhoto", "Failed to update profile: %v", err)
				return nil, err
			}
		}

		reloaded, err := s.repo.User.FindByID(tx, userID, "Roles", "Profile")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ProfileUploadFacePhoto", "Failed to upload face photo: %v", err)
		return nil, err
	}

	result := res.(*models.User)
	dto := dtos.ToProfileDTO(result)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "Face Photo Uploaded",
		Message: "Your face photo has been uploaded successfully",
		Data: map[string]interface{}{
			"user_id": userID,
		},
	})

	s.Logger.LogEnd("ProfileUploadFacePhoto", "Face photo uploaded for user: %s", dto.Email)
	return &dto, nil
}

// ProfileDeleteFacePhoto removes the face photo for the authenticated user.
func (s *Services) ProfileDeleteFacePhoto(ctx context.Context, userID uint) (*dtos.ProfileDTO, error) {
	s.Logger.LogStart("ProfileDeleteFacePhoto", "Deleting face photo for user ID: %d", userID)

	existing, err := s.repo.User.FindByID(nil, userID, "Profile")
	if err != nil {
		s.Logger.LogEndWithError("ProfileDeleteFacePhoto", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	if existing.Profile == nil || existing.Profile.FacePhotoURL == "" {
		s.Logger.LogEndWithError("ProfileDeleteFacePhoto", "No face photo to delete")
		return nil, &helpers.FieldError{Field: "face_photo", Message: "No face photo to delete"}
	}

	oldFacePhoto := existing.Profile.FacePhotoURL

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		profileUpdates := map[string]interface{}{
			"face_photo_url":  "",
			"face_embedding": "",
		}

		if _, err := s.repo.UserProfile.UpdateMap(tx, &models.UserProfile{UserID: userID}, profileUpdates); err != nil {
			s.Logger.LogStep("ProfileDeleteFacePhoto", "Failed to update profile: %v", err)
			return nil, err
		}

		reloaded, err := s.repo.User.FindByID(tx, userID, "Roles", "Profile")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ProfileDeleteFacePhoto", "Failed to delete face photo: %v", err)
		return nil, err
	}

	// Delete the file after successful DB update
	helpers.DeleteFile(oldFacePhoto)

	result := res.(*models.User)
	dto := dtos.ToProfileDTO(result)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "Face Photo Removed",
		Message: "Your face photo has been removed",
		Data: map[string]interface{}{
			"user_id": userID,
		},
	})

	s.Logger.LogEnd("ProfileDeleteFacePhoto", "Face photo deleted for user: %s", dto.Email)
	return &dto, nil
}
