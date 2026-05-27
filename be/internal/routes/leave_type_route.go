package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

func RegisterLeaveTypeRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	leaveTypes := r.Group("/leave/types")
	{
		leaveTypes.POST("", middleware.RequirePermission(acc, "leave.manage-types"), handlers.LeaveTypeCreate)
		leaveTypes.GET("", middleware.RequirePermission(acc, "leave.manage-types"), handlers.LeaveTypeGetAll)
		leaveTypes.GET("/:id", middleware.RequirePermission(acc, "leave.manage-types"), handlers.LeaveTypeGetByID)
		leaveTypes.PUT("/:id", middleware.RequirePermission(acc, "leave.manage-types"), handlers.LeaveTypeUpdate)
		leaveTypes.DELETE("/:id", middleware.RequirePermission(acc, "leave.manage-types"), handlers.LeaveTypeDelete)
	}
}
