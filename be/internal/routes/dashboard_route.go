package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

func RegisterDashboardRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/employee", handlers.EmployeeDashboard)
		dashboard.GET("/hr", middleware.RequirePermission(acc, "dashboard.view-hr"), handlers.HRDashboard)
	}
}

func RegisterScheduleRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	r.GET("/shifts/schedule", middleware.RequirePermission(acc, "shift-assign.index"), handlers.AttendanceSchedule)
}
