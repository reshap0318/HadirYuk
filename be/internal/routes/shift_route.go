package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

func RegisterShiftRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	shifts := r.Group("/shifts")
	{
		shifts.POST("", middleware.RequirePermission(acc, "shift.create"), handlers.ShiftCreate)
		shifts.GET("", middleware.RequirePermission(acc, "shift.index"), handlers.ShiftGetAll)
		shifts.GET("/:id", middleware.RequirePermission(acc, "shift.index"), handlers.ShiftGetByID)
		shifts.PUT("/:id", middleware.RequirePermission(acc, "shift.update"), handlers.ShiftUpdate)
		shifts.DELETE("/:id", middleware.RequirePermission(acc, "shift.delete"), handlers.ShiftDelete)
	}
}
