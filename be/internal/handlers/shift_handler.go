package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

func (h *Handlers) ShiftCreate(c *gin.Context) {
	var req dtos.ShiftRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.ShiftCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create shift") {
		return
	}

	helpers.Created(c, "Shift created successfully", dto)
}

func (h *Handlers) ShiftGetAll(c *gin.Context) {
	pageStr := c.Query("page")

	if pageStr == "" {
		shifts, err := h.svcs.ShiftGetAllUnpaginated(c.Request.Context())
		if err != nil {
			helpers.InternalServerError(c, "Failed to fetch shifts")
			return
		}

		helpers.OK(c, "Shifts fetched successfully", shifts)
		return
	}

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.ShiftGetAllPaginated(c.Request.Context(), opts)
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch shifts")
		return
	}

	helpers.OKWithMetadata(c, "Shifts fetched successfully", result)
}

func (h *Handlers) ShiftGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid shift ID")
		return
	}

	dto, err := h.svcs.ShiftGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch shift") {
		return
	}

	helpers.OK(c, "Shift fetched successfully", dto)
}

func (h *Handlers) ShiftUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid shift ID")
		return
	}

	var req dtos.ShiftRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.ShiftUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update shift") {
		return
	}

	helpers.OK(c, "Shift updated successfully", dto)
}

func (h *Handlers) ShiftDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid shift ID")
		return
	}

	err = h.svcs.ShiftDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete shift") {
		return
	}

	helpers.OK(c, "Shift deleted successfully", nil)
}

func (h *Handlers) ShiftAssignToUser(c *gin.Context) {
	var req dtos.ShiftAssignmentRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.ShiftAssignToUser(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to assign shift") {
		return
	}

	helpers.Created(c, "Shift assigned successfully", dto)
}

func (h *Handlers) ShiftGetUserAssignments(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		userIDStr = c.Query("user_id")
	}

	if userIDStr == "" {
		// List all assignments (admin view)
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		search := c.Query("search")

		opts := &repositories.QueryOptions{
			Page:         page,
			PageSize:     pageSize,
			Preloads:     []string{"User", "Shift"},
			Search:       search,
			SearchFields: []string{"users.name", "users.email", "shifts.name"},
		}

		result, err := h.svcs.ShiftGetAllAssignments(c.Request.Context(), opts)
		if err != nil {
			helpers.InternalServerError(c, "Failed to fetch assignments")
			return
		}

		helpers.OKWithMetadata(c, "Assignments fetched successfully", result)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.ShiftGetUserAssignments(c.Request.Context(), uint(userID), opts)
	if helpers.HandleError(c, err, "Failed to fetch assignments") {
		return
	}

	helpers.OKWithMetadata(c, "Assignments fetched successfully", result)
}

func (h *Handlers) ShiftGetUserActiveAssignment(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	dto, err := h.svcs.ShiftGetUserActiveAssignment(c.Request.Context(), uint(userID))
	if err == helpers.ErrNotFound {
		helpers.OK(c, "No active assignment found", nil)
		return
	}
	if helpers.HandleError(c, err, "Failed to fetch active assignment") {
		return
	}

	helpers.OK(c, "Active assignment fetched successfully", dto)
}

func (h *Handlers) ShiftUpdateAssignment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid assignment ID")
		return
	}

	var req dtos.ShiftAssignmentUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.ShiftUpdateAssignment(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update assignment") {
		return
	}

	helpers.OK(c, "Assignment updated successfully", dto)
}

func (h *Handlers) ShiftDeleteAssignment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid assignment ID")
		return
	}

	err = h.svcs.ShiftDeleteAssignment(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete assignment") {
		return
	}

	helpers.OK(c, "Assignment deleted successfully", nil)
}
