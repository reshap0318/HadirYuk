package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

// RegisterPermissionRoutes registers protected permission routes.
func RegisterPermissionRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	permissions := r.Group("/permissions")
	{
		permissions.POST("", middleware.RequirePermission(acc, "permission.create"), handlers.PermissionCreate)
		permissions.GET("", middleware.RequirePermission(acc, "permission.index"), handlers.PermissionGetAll)
		permissions.GET("/:id", middleware.RequirePermission(acc, "permission.index"), handlers.PermissionGetByID)
		permissions.PUT("/:id", middleware.RequirePermission(acc, "permission.edit"), handlers.PermissionUpdate)
		permissions.DELETE("/:id", middleware.RequirePermission(acc, "permission.delete"), handlers.PermissionDelete)
	}
}
