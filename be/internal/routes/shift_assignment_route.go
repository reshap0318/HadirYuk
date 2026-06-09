package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

func RegisterShiftAssignmentRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	shifts := r.Group("/shifts")
	{
		// Assignment routes
		shifts.POST("/assignments", middleware.RequirePermission(acc, "shift-assign.create"), handlers.ShiftAssignToUser)
		shifts.GET("/assignments", middleware.RequirePermission(acc, "shift-assign.index"), handlers.ShiftGetUserAssignments)
		shifts.GET("/assignments/:user_id", middleware.RequirePermission(acc, "shift-assign.index"), handlers.ShiftGetUserAssignments)
		shifts.GET("/assignments/:user_id/active", middleware.RequirePermission(acc, "shift-assign.index"), handlers.ShiftGetUserActiveAssignment)
		shifts.PUT("/assignments/:id", middleware.RequirePermission(acc, "shift-assign.update"), handlers.ShiftUpdateAssignment)
		shifts.DELETE("/assignments/:id", middleware.RequirePermission(acc, "shift-assign.delete"), handlers.ShiftDeleteAssignment)
	}
}
