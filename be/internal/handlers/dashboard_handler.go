package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/helpers"
)

func (h *Handlers) EmployeeDashboard(c *gin.Context) {
	data, err := h.svcs.EmployeeDashboard(c.Request.Context())
	if helpers.HandleError(c, err, "Failed to fetch employee dashboard") {
		return
	}

	helpers.OK(c, "Employee dashboard data fetched successfully", data)
}

func (h *Handlers) HRDashboard(c *gin.Context) {
	data, err := h.svcs.HRDashboard(c.Request.Context())
	if helpers.HandleError(c, err, "Failed to fetch HR dashboard") {
		return
	}

	helpers.OK(c, "HR dashboard data fetched successfully", data)
}

func (h *Handlers) AttendanceSchedule(c *gin.Context) {
	dateFrom := c.DefaultQuery("date_from", time.Now().Format("2006-01-02"))
	dateTo := c.DefaultQuery("date_to", time.Now().AddDate(0, 0, 7).Format("2006-01-02"))

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "50")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	data, err := h.svcs.AttendanceSchedule(c.Request.Context(), dateFrom, dateTo)
	if helpers.HandleError(c, err, "Failed to fetch attendance schedule") {
		return
	}

	// Manual pagination
	total := len(data)
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	var pagedData interface{}
	if start >= total {
		pagedData = []interface{}{}
	} else if end > total {
		pagedData = data[start:]
	} else {
		pagedData = data[start:end]
	}

	helpers.OKWithMetadata(c, "Schedule fetched successfully", &schedulePaginatedResult{
		data:       pagedData,
		total:      int64(total),
		page:       page,
		pageSize:   pageSize,
		totalPages: totalPages,
	})
}

type schedulePaginatedResult struct {
	data       interface{}
	total      int64
	page       int
	pageSize   int
	totalPages int
}

func (r *schedulePaginatedResult) GetData() interface{} {
	return r.data
}

func (r *schedulePaginatedResult) GetMetadata() helpers.PaginationMeta {
	return helpers.PaginationMeta{
		Total:      r.total,
		Page:       r.page,
		PageSize:   r.pageSize,
		TotalPages: r.totalPages,
	}
}
