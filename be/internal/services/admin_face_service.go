package services

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/clients/face"
	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
)

// AdminFacePhotoUpload uploads a face photo for a specific user (admin action).
func (s *Services) AdminFacePhotoUpload(ctx context.Context, userID uint, fileUUID string) (*dtos.FacePhotoResponse, error) {
	s.Logger.LogStart("AdminFacePhotoUpload", "Uploading face photo for user ID: %d", userID)

	// Find the target user
	user, err := s.repo.User.FindByID(nil, userID, "Profile")
	if err != nil {
		s.Logger.LogEndWithError("AdminFacePhotoUpload", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	// Move file from tmp to face-photos directory first
	// Embedding will be computed from the final path to ensure consistency
	facePhotoPath, err := helpers.MoveFile(fileUUID, "storage/tmp", "storage/face-photos")
	if err != nil {
		s.Logger.LogEndWithError("AdminFacePhotoUpload", "Failed to move face photo: %v", err)
		return nil, err
	}

	// Process face photo to generate embedding from final path
	// This ensures embedding is consistent for face matching (not path-dependent)
	embedding, err := s.FaceService.ProcessFacePhoto(facePhotoPath)
	if err != nil {
		s.Logger.LogStep("AdminFacePhotoUpload", "Embedding generation failed: %v", err)
	}

	// Serialize embedding to JSON
	embeddingJSON := ""
	if embedding != nil {
		embeddingBytes, marshalErr := json.Marshal(embedding)
		if marshalErr != nil {
			s.Logger.LogStep("AdminFacePhotoUpload", "Failed to marshal embedding: %v", marshalErr)
		} else {
			embeddingJSON = string(embeddingBytes)
		}
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		// Delete old face photo if exists
		if user.Profile != nil && user.Profile.FacePhotoURL != "" {
			helpers.DeleteFile(user.Profile.FacePhotoURL)
		}

		profileUpdates := map[string]interface{}{
			"face_photo_url": facePhotoPath,
		}
		if embeddingJSON != "" {
			profileUpdates["face_embedding"] = embeddingJSON
		}

		if user.Profile == nil {
			profile := &models.UserProfile{
				UserID:        userID,
				FacePhotoURL:  facePhotoPath,
				FaceEmbedding: embeddingJSON,
			}
			if _, err := s.repo.UserProfile.Create(tx, profile); err != nil {
				s.Logger.LogStep("AdminFacePhotoUpload", "Failed to create profile: %v", err)
				return nil, err
			}
		} else {
			if _, err := s.repo.UserProfile.UpdateMap(tx, &models.UserProfile{UserID: userID}, profileUpdates); err != nil {
				s.Logger.LogStep("AdminFacePhotoUpload", "Failed to update profile: %v", err)
				return nil, err
			}
		}

		return nil, nil
	})
	if err != nil {
		// Clean up moved file on transaction failure
		helpers.DeleteFile(facePhotoPath)
		s.Logger.LogEndWithError("AdminFacePhotoUpload", "Failed to upload face photo: %v", err)
		return nil, err
	}

	_ = res

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "Face Photo Uploaded (Admin)",
		Message: fmt.Sprintf("Face photo uploaded for user ID: %d", userID),
		Data: map[string]interface{}{
			"user_id": userID,
		},
	})

	response := &dtos.FacePhotoResponse{
		PhotoURL:       helpers.GetFileURL(facePhotoPath),
		EmbeddingReady: embeddingJSON != "",
	}

	s.Logger.LogEnd("AdminFacePhotoUpload", "Face photo uploaded for user ID: %d", userID)
	return response, nil
}

// AdminFacePhotoDelete removes face photo for a specific user (admin action).
func (s *Services) AdminFacePhotoDelete(ctx context.Context, userID uint) (*dtos.FacePhotoResponse, error) {
	s.Logger.LogStart("AdminFacePhotoDelete", "Deleting face photo for user ID: %d", userID)

	user, err := s.repo.User.FindByID(nil, userID, "Profile")
	if err != nil {
		s.Logger.LogEndWithError("AdminFacePhotoDelete", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	if user.Profile == nil || user.Profile.FacePhotoURL == "" {
		s.Logger.LogEndWithError("AdminFacePhotoDelete", "No face photo to delete")
		return nil, &helpers.FieldError{Field: "face_photo", Message: "No face photo to delete"}
	}

	oldFacePhoto := user.Profile.FacePhotoURL

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		profileUpdates := map[string]interface{}{
			"face_photo_url":  "",
			"face_embedding": "",
		}

		if _, err := s.repo.UserProfile.UpdateMap(tx, &models.UserProfile{UserID: userID}, profileUpdates); err != nil {
			s.Logger.LogStep("AdminFacePhotoDelete", "Failed to update profile: %v", err)
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("AdminFacePhotoDelete", "Failed to delete face photo: %v", err)
		return nil, err
	}

	_ = res

	// Delete the file after successful DB update
	helpers.DeleteFile(oldFacePhoto)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "Face Photo Removed (Admin)",
		Message: fmt.Sprintf("Face photo removed for user ID: %d", userID),
		Data: map[string]interface{}{
			"user_id": userID,
		},
	})

	response := &dtos.FacePhotoResponse{
		PhotoURL:       "",
		EmbeddingReady: false,
	}

	s.Logger.LogEnd("AdminFacePhotoDelete", "Face photo deleted for user ID: %d", userID)
	return response, nil
}

// FaceMatch performs face matching against all users with face embeddings.
// Phase 1 stub: returns mock match result.
// Phase 2: decode base64 → temp file → detect → embed → match against all users.
func (s *Services) FaceMatch(ctx context.Context, photoBase64 string) (*dtos.FaceMatchResponse, error) {
	s.Logger.LogStart("FaceMatch", "Performing face match")

	// Decode base64 to temp file
	tempFile, err := face.Base64ToTempFile(photoBase64)
	if err != nil {
		s.Logger.LogEndWithError("FaceMatch", "Failed to decode base64 image: %v", err)
		return nil, &helpers.FieldError{Field: "photo", Message: "Invalid image format or size"}
	}
	defer face.CleanupTempFile(tempFile)

	// Validate image format
	if err := face.ValidateFaceImageFormat(tempFile); err != nil {
		s.Logger.LogEndWithError("FaceMatch", "Invalid image format: %v", err)
		return nil, &helpers.FieldError{Field: "photo", Message: "Format file harus JPG/PNG/WebP"}
	}

	// Process face photo to generate embedding
	embedding, err := s.FaceService.ProcessFacePhoto(tempFile)
	if err != nil {
		s.Logger.LogEndWithError("FaceMatch", "Face processing failed: %v", err)
		return nil, &helpers.FieldError{Field: "photo", Message: "Wajah tidak terdeteksi. Pastikan wajah terlihat jelas"}
	}

	// Find users with face embeddings only (optimization: skip users without embeddings)
	users, err := s.repo.User.FindByCondition(nil, "face_embedding IS NOT NULL AND face_embedding != ''", "Profile")
	if err != nil {
		s.Logger.LogEndWithError("FaceMatch", "Failed to fetch users: %v", err)
		return nil, err
	}

	// Match against each user's embedding
	var bestMatch *dtos.FaceMatchResponse
	bestSimilarity := 0.0

	for _, u := range users {
		if u.Profile == nil || u.Profile.FaceEmbedding == "" {
			continue
		}

		var storedEmbedding []float64
		if err := json.Unmarshal([]byte(u.Profile.FaceEmbedding), &storedEmbedding); err != nil {
			s.Logger.LogStep("FaceMatch", "Failed to unmarshal embedding for user %d: %v", u.ID, err)
			continue
		}

		similarity, matched := s.FaceService.MatchFace(embedding, storedEmbedding)
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestMatch = &dtos.FaceMatchResponse{
				Matched:    matched,
				Similarity: similarity,
				UserID:     u.ID,
			}
		}
	}

	if bestMatch == nil {
		s.Logger.LogEnd("FaceMatch", "No users with face embeddings found")
		return &dtos.FaceMatchResponse{
			Matched:    false,
			Similarity: 0,
		}, nil
	}

	if !bestMatch.Matched {
		s.Logger.LogEnd("FaceMatch", "Face not recognized, similarity: %.4f", bestMatch.Similarity)
		return bestMatch, nil
	}

	s.Logger.LogEnd("FaceMatch", "Face matched with user ID: %d, similarity: %.4f", bestMatch.UserID, bestMatch.Similarity)
	return bestMatch, nil
}
