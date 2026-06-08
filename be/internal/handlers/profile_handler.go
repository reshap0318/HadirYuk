package handlers

import (
	"github.com/gin-gonic/gin"

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
