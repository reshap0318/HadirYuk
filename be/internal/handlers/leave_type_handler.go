package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

func (h *Handlers) LeaveTypeCreate(c *gin.Context) {
	var req dtos.LeaveTypeRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.LeaveTypeCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create leave type") {
		return
	}

	helpers.Created(c, "Leave type created successfully", dto)
}

func (h *Handlers) LeaveTypeGetAll(c *gin.Context) {
	pageStr := c.Query("page")

	if pageStr == "" {
		leaveTypes, err := h.svcs.LeaveTypeGetAllUnpaginated(c.Request.Context())
		if err != nil {
			helpers.InternalServerError(c, "Failed to fetch leave types")
			return
		}

		helpers.OK(c, "Leave types fetched successfully", leaveTypes)
		return
	}

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.LeaveTypeGetAllPaginated(c.Request.Context(), opts)
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch leave types")
		return
	}

	helpers.OKWithMetadata(c, "Leave types fetched successfully", result)
}

func (h *Handlers) LeaveTypeGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid leave type ID")
		return
	}

	dto, err := h.svcs.LeaveTypeGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch leave type") {
		return
	}

	helpers.OK(c, "Leave type fetched successfully", dto)
}

func (h *Handlers) LeaveTypeUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid leave type ID")
		return
	}

	var req dtos.LeaveTypeRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.LeaveTypeUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update leave type") {
		return
	}

	helpers.OK(c, "Leave type updated successfully", dto)
}

func (h *Handlers) LeaveTypeDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid leave type ID")
		return
	}

	err = h.svcs.LeaveTypeDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete leave type") {
		return
	}

	helpers.OK(c, "Leave type deleted successfully", nil)
}
