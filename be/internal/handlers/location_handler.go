package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

func (h *Handlers) LocationCreate(c *gin.Context) {
	var req dtos.LocationRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.LocationCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create location") {
		return
	}

	helpers.Created(c, "Location created successfully", dto)
}

func (h *Handlers) LocationGetAll(c *gin.Context) {
	pageStr := c.Query("page")

	if pageStr == "" {
		locations, err := h.svcs.LocationGetAllUnpaginated(c.Request.Context())
		if err != nil {
			helpers.InternalServerError(c, "Failed to fetch locations")
			return
		}

		helpers.OK(c, "Locations fetched successfully", locations)
		return
	}

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.LocationGetAllPaginated(c.Request.Context(), opts)
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch locations")
		return
	}

	helpers.OKWithMetadata(c, "Locations fetched successfully", result)
}

func (h *Handlers) LocationGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid location ID")
		return
	}

	dto, err := h.svcs.LocationGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch location") {
		return
	}

	helpers.OK(c, "Location fetched successfully", dto)
}

func (h *Handlers) LocationUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid location ID")
		return
	}

	var req dtos.LocationRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.LocationUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update location") {
		return
	}

	helpers.OK(c, "Location updated successfully", dto)
}

func (h *Handlers) LocationDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid location ID")
		return
	}

	err = h.svcs.LocationDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete location") {
		return
	}

	helpers.OK(c, "Location deleted successfully", nil)
}
