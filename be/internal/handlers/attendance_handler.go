package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
)

func (h *Handlers) GetTodayStatus(c *gin.Context) {
	status, err := h.svcs.GetTodayStatus(c.Request.Context())
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.OK(c, "Today's attendance status", status)
}

func (h *Handlers) AttendanceCheckIn(c *gin.Context) {
	var req dtos.AttendanceCheckInRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.AttendanceCheckIn(c.Request.Context(), req)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.Created(c, "Check-in berhasil", dto)
}

func (h *Handlers) AttendanceCheckOut(c *gin.Context) {
	var req dtos.AttendanceCheckInRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.AttendanceCheckOut(c.Request.Context(), req)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.OK(c, "Check-out berhasil", dto)
}

func (h *Handlers) GetNearestOffice(c *gin.Context) {
	var req dtos.NearestOfficeRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	response, err := h.svcs.NearestOffice(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to find nearest office") {
		return
	}

	helpers.OK(c, "Nearest office found", response)
}
