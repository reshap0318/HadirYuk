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
