package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
)

// AdminFacePhotoUpload handles POST /api/users/:id/face-photo
func (h *Handlers) AdminFacePhotoUpload(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	var req dtos.FacePhotoRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Face photo UUID is required")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	response, err := h.svcs.AdminFacePhotoUpload(c.Request.Context(), uint(userID), req.FacePhoto)
	if helpers.HandleError(c, err, "Failed to upload face photo") {
		return
	}

	helpers.OK(c, "Face photo uploaded successfully", response)
}

// AdminFacePhotoDelete handles DELETE /api/users/:id/face-photo
func (h *Handlers) AdminFacePhotoDelete(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	response, err := h.svcs.AdminFacePhotoDelete(c.Request.Context(), uint(userID))
	if helpers.HandleError(c, err, "Failed to delete face photo") {
		return
	}

	helpers.OK(c, "Face photo deleted successfully", response)
}

// FaceMatch handles POST /api/face/match (for attendance verification)
func (h *Handlers) FaceMatch(c *gin.Context) {
	var req dtos.FaceMatchRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Photo base64 is required")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	response, err := h.svcs.FaceMatch(c.Request.Context(), req.PhotoBase64)
	if helpers.HandleError(c, err, "Failed to perform face match") {
		return
	}

	helpers.OK(c, "Face match completed", response)
}
