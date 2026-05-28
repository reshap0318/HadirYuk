package services

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

// UserCreate creates a new user with optional roles.
func (s *Services) UserCreate(ctx context.Context, req dtos.UserCreateRequest) (*dtos.UserDTO, error) {
	s.Logger.LogStart("UserCreate", "Creating user: %s", req.Email)

	exists, err := s.repo.User.Exists(nil, map[string]interface{}{"email": req.Email})
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to check email: %v", err)
		return nil, err
	}
	if exists {
		s.Logger.LogEndWithError("UserCreate", "Email already exists: %s", req.Email)
		return nil, &helpers.FieldError{Field: "email", Message: "user already exists"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to hash password: %v", err)
		return nil, err
	}

	avatarPath := ""
	if req.Avatar != "" {
		avatarPath, err = helpers.MoveFile(req.Avatar, "storage/tmp", "storage/avatars")
		if err != nil {
			s.Logger.LogStep("UserCreate", "Failed to move avatar: %v", err)
			avatarPath = ""
		}
	}

	user := &models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
		Avatar:   avatarPath,
		Profile: &models.UserProfile{
			Department: req.Department,
			Position:   req.Position,
			Phone:      req.Phone,
			JoinDate:   req.JoinDate,
		},
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		result, err := s.repo.User.Create(tx, user)
		if err != nil {
			return nil, err
		}

		var roles []models.Role
		for _, roleID := range req.Roles {
			roles = append(roles, models.Role{ID: roleID})
		}
		if err := tx.Model(&result).Association("Roles").Append(roles); err != nil {
			s.Logger.LogStep("UserCreate", "Failed to assign roles: %v", err)
		}

		reloaded, err := s.repo.User.FindByID(tx, result.ID, "Roles", "Profile")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to create user: %v", err)
		return nil, err
	}

	result := res.(*models.User)
	dto := dtos.ToUserDTO(result)
	s.Logger.LogEnd("UserCreate", "User created: %s (ID: %d)", dto.Email, dto.ID)
	return &dto, nil
}

// UserGetAll returns all users with roles and profiles (no pagination).
func (s *Services) UserGetAll(ctx context.Context) ([]dtos.UserDTO, error) {
	users, err := s.repo.User.FindAll(nil, "Roles", "Profile")
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dtos.UserDTO, len(users))
	for i, u := range users {
		userDTOs[i] = dtos.ToUserDTO(&u)
	}
	return userDTOs, nil
}

// UserGetAllPaginated returns paginated users with roles.
func (s *Services) UserGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.UserDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "ASC"
	}
	opts.Preloads = []string{"Roles", "Profile"}

	result, err := s.repo.User.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dtos.UserDTO, len(result.Data))
	for i, u := range result.Data {
		userDTOs[i] = dtos.ToUserDTO(&u)
	}

	return &repositories.PagedResult[dtos.UserDTO]{
		Data:       userDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// UserGetByID returns a user by ID with roles.
func (s *Services) UserGetByID(ctx context.Context, id uint) (*dtos.UserDTO, error) {
	user, err := s.repo.User.FindByID(nil, id, "Roles", "Profile")
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToUserDTO(user)
	return &dto, nil
}

// ProfileGet returns the authenticated user's profile.
func (s *Services) ProfileGet(ctx context.Context, userID uint) (*dtos.UserDTO, error) {
	s.Logger.LogStart("ProfileGet", "Fetching profile for user ID: %d", userID)

	user, err := s.repo.User.FindByID(nil, userID, "Roles", "Profile")
	if err != nil {
		s.Logger.LogEndWithError("ProfileGet", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToUserDTO(user)
	s.Logger.LogEnd("ProfileGet", "Profile fetched for user: %s", dto.Email)
	return &dto, nil
}

// ProfileUpdate updates the authenticated user's profile.
func (s *Services) ProfileUpdate(ctx context.Context, userID uint, req dtos.ProfileUpdateRequest) (*dtos.UserDTO, error) {
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
	dto := dtos.ToUserDTO(result)

	s.Access.Invalidate(userID)

	s.Logger.LogEnd("ProfileUpdate", "Profile updated for user: %s", dto.Email)
	return &dto, nil
}

// UserUpdate updates an existing user with optional roles and profile.
func (s *Services) UserUpdate(ctx context.Context, id uint, req dtos.UserUpdateRequest) (*dtos.UserDTO, error) {
	s.Logger.LogStart("UserUpdate", "Updating user ID: %d", id)

	existing, err := s.repo.User.FindByID(nil, id, "Profile")
	if err != nil {
		s.Logger.LogEndWithError("UserUpdate", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	if existing.Email != req.Email {
		exists, err := s.repo.User.Exists(nil, map[string]interface{}{"email": req.Email})
		if err != nil {
			s.Logger.LogEndWithError("UserUpdate", "Failed to check email: %v", err)
			return nil, err
		}
		if exists {
			s.Logger.LogEndWithError("UserUpdate", "Email already exists: %s", req.Email)
			return nil, &helpers.FieldError{Field: "email", Message: "user already exists"}
		}
	}

	updates := map[string]interface{}{
		"name":  req.Name,
		"email": req.Email,
	}

	if req.Avatar != "" {
		avatarPath, err := helpers.MoveFile(req.Avatar, "storage/tmp", "storage/avatars")
		if err != nil {
			s.Logger.LogStep("UserUpdate", "Failed to move avatar: %v", err)
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
		result, err := s.repo.User.UpdateMap(tx, &models.User{ID: id}, updates)
		if err != nil {
			return nil, err
		}

		if len(profileUpdates) > 0 {
			if result.Profile == nil {
				result.Profile = &models.UserProfile{UserID: id}
				if _, err := s.repo.UserProfile.Create(tx, result.Profile); err != nil {
					s.Logger.LogStep("UserUpdate", "Failed to create profile: %v", err)
					return nil, err
				}
			}
			if _, err := s.repo.UserProfile.UpdateMap(tx, &models.UserProfile{UserID: id}, profileUpdates); err != nil {
				s.Logger.LogStep("UserUpdate", "Failed to update profile fields: %v", err)
				return nil, err
			}
		}

		if err := tx.Model(&result).Association("Roles").Clear(); err != nil {
			return nil, err
		}

		var roles []models.Role
		for _, roleID := range req.Roles {
			roles = append(roles, models.Role{ID: roleID})
		}
		if err := tx.Model(&result).Association("Roles").Append(roles); err != nil {
			s.Logger.LogStep("UserUpdate", "Failed to assign roles: %v", err)
			return nil, err
		}

		reloaded, err := s.repo.User.FindByID(tx, result.ID, "Roles", "Profile")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("UserUpdate", "Failed to update user: %v", err)
		return nil, err
	}

	result := res.(*models.User)
	dto := dtos.ToUserDTO(result)

	// Invalidate cached session so next request gets updated permissions
	s.Access.Invalidate(id)

	s.Logger.LogEnd("UserUpdate", "User updated: %s (ID: %d)", dto.Email, dto.ID)
	return &dto, nil
}

// UserDelete soft deletes a user and its role associations.
func (s *Services) UserDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("UserDelete", "Deleting user ID: %d", id)

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		user := models.User{ID: id}
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
		_, err := s.repo.User.Delete(tx, id)
		return err
	}); err != nil {
		s.Logger.LogEndWithError("UserDelete", "Failed to delete user: %v", err)
		return err
	}

	// Invalidate cached session
	s.Access.Invalidate(id)

	s.Logger.LogEnd("UserDelete", "User deleted: ID: %d", id)
	return nil
}

// ProfileUploadFacePhoto uploads a face photo for the authenticated user.
func (s *Services) ProfileUploadFacePhoto(ctx context.Context, userID uint, fileUUID string) (*dtos.UserDTO, error) {
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
	dto := dtos.ToUserDTO(result)

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
func (s *Services) ProfileDeleteFacePhoto(ctx context.Context, userID uint) (*dtos.UserDTO, error) {
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
	dto := dtos.ToUserDTO(result)

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
