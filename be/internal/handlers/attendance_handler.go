package handlers

import (
	"fmt"
	"strconv"

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

func (h *Handlers) AttendanceQRCheckIn(c *gin.Context) {
	var req dtos.AttendanceQRCheckInRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.AttendanceQRCheckIn(c.Request.Context(), req)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.Created(c, "Check-in berhasil", dto)
}

func (h *Handlers) AttendanceQRCheckOut(c *gin.Context) {
	var req dtos.AttendanceQRCheckInRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.AttendanceQRCheckOut(c.Request.Context(), req)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.OK(c, "Check-out berhasil", dto)
}

func (h *Handlers) AttendanceHistory(c *gin.Context) {
	var req dtos.AttendanceHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	// Check if HR wants all users
	userParam := c.Query("user")
	if userParam == "all" {
		result, err := h.svcs.AttendanceHistoryAll(c.Request.Context(), req)
		if helpers.HandleError(c, err, "") {
			return
		}
		helpers.OKWithMetadata(c, "Attendance history retrieved", result)
		return
	}

	result, err := h.svcs.AttendanceHistory(c.Request.Context(), req)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.OKWithMetadata(c, "Attendance history retrieved", result)
}

func (h *Handlers) AttendanceMonthlyStats(c *gin.Context) {
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")

	if year == "" || month == "" {
		helpers.BadRequest(c, "Year and month are required")
		return
	}

	var y, m int
	if _, err := fmt.Sscanf(year, "%d", &y); err != nil {
		helpers.BadRequest(c, "Invalid year format")
		return
	}
	if _, err := fmt.Sscanf(month, "%d", &m); err != nil || m < 1 || m > 12 {
		helpers.BadRequest(c, "Invalid month format")
		return
	}

	stats, err := h.svcs.AttendanceMonthlyStats(c.Request.Context(), y, m)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.OK(c, "Monthly statistics retrieved", stats)
}

func (h *Handlers) AttendanceCorrect(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		helpers.BadRequest(c, "Invalid attendance ID")
		return
	}

	var req dtos.AttendanceCorrectRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.AttendanceCorrect(c.Request.Context(), id, req)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.OK(c, "Attendance record corrected", dto)
}

func (h *Handlers) AttendanceReport(c *gin.Context) {
	var req dtos.AttendanceReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.svcs.AttendanceReport(c.Request.Context(), req)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.OKWithMetadata(c, "Attendance report retrieved", result)
}

func (h *Handlers) AttendanceLateStats(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	userIDStr := c.Query("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var userID *uint
	if userIDStr != "" {
		var id uint
		if _, err := fmt.Sscanf(userIDStr, "%d", &id); err == nil {
			userID = &id
		}
	}

	stats, err := h.svcs.AttendanceLateStats(c.Request.Context(), &dateFrom, &dateTo, userID, page, pageSize)
	if helpers.HandleError(c, err, "") {
		return
	}

	helpers.OK(c, "Late statistics retrieved", stats)
}
