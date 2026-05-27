package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

// RegisterRoleRoutes registers protected role routes.
func RegisterRoleRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	roles := r.Group("/roles")
	{
		roles.POST("", middleware.RequirePermission(acc, "role.create"), handlers.RoleCreate)
		roles.GET("", middleware.RequirePermission(acc, "role.index"), handlers.RoleGetAll)
		roles.GET("/:id", middleware.RequirePermission(acc, "role.index"), handlers.RoleGetByID)
		roles.PUT("/:id", middleware.RequirePermission(acc, "role.edit"), handlers.RoleUpdate)
		roles.DELETE("/:id", middleware.RequirePermission(acc, "role.delete"), handlers.RoleDelete)
		roles.GET("/:id/permissions", middleware.RequirePermission(acc, "role.index"), handlers.RoleGetPermissions)
	}
}
