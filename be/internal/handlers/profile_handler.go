package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
)

// ProfileGet handles GET /api/me
func (h *Handlers) ProfileGet(c *gin.Context) {
	userID := c.GetUint("user_id")

	dto, err := h.svcs.ProfileGet(c.Request.Context(), userID)
	if helpers.HandleError(c, err, "Failed to fetch profile") {
		return
	}

	helpers.OK(c, "Profile fetched successfully", dto)
}

// ProfileUpdate handles PUT /api/me
func (h *Handlers) ProfileUpdate(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dtos.ProfileUpdateRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.ProfileUpdate(c.Request.Context(), userID, req)
	if helpers.HandleError(c, err, "Failed to update profile") {
		return
	}

	helpers.OK(c, "Profile updated successfully", dto)
}

// ProfileUploadFacePhoto handles POST /api/me/face-photo
func (h *Handlers) ProfileUploadFacePhoto(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, header, err := c.Request.FormFile("face_photo")
	if err != nil {
		helpers.BadRequest(c, "Face photo file is required")
		return
	}
	defer file.Close()

	// Use the existing file upload logic
	fileUUID, filePath, err := h.uploadFacePhotoFile(c, header.Filename)
	if err != nil {
		helpers.InternalServerError(c, fmt.Sprintf("Failed to upload face photo: %v", err))
		return
	}

	dto, err := h.svcs.ProfileUploadFacePhoto(c.Request.Context(), userID, fileUUID)
	if helpers.HandleError(c, err, "Failed to upload face photo") {
		// Clean up the uploaded file if service fails
		helpers.DeleteFile(filePath)
		return
	}

	helpers.OK(c, "Face photo uploaded successfully", dto)
}

// ProfileDeleteFacePhoto handles DELETE /api/me/face-photo
func (h *Handlers) ProfileDeleteFacePhoto(c *gin.Context) {
	userID := c.GetUint("user_id")

	dto, err := h.svcs.ProfileDeleteFacePhoto(c.Request.Context(), userID)
	if helpers.HandleError(c, err, "Failed to delete face photo") {
		return
	}

	helpers.OK(c, "Face photo deleted successfully", dto)
}

// uploadFacePhotoFile handles the file upload for face photos
func (h *Handlers) uploadFacePhotoFile(c *gin.Context, originalFilename string) (string, string, error) {
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext == "" {
		return "", "", fmt.Errorf("file extension is required")
	}

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return "", "", fmt.Errorf("file type %s is not allowed, only jpg, jpeg, png, webp are accepted", ext)
	}

	fileUUID := uuid.New().String()
	fileName := fmt.Sprintf("%s%s", fileUUID, ext)
	uploadDir := "storage/tmp"

	filePath, err := helpers.SaveUploadedFileWithOpts(c, "face_photo", uploadDir, &helpers.SaveFileOptions{
		CustomName: fileName,
	})
	if err != nil {
		return "", "", err
	}

	return fileUUID, filePath, nil
}
