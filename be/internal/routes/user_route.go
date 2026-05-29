package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

// RegisterUserRoutes registers protected user routes.
func RegisterUserRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	users := r.Group("/users")
	{
		users.POST("", middleware.RequirePermission(acc, "user.create"), handlers.UserCreate)
		users.GET("", middleware.RequirePermission(acc, "user.index", "user.view-all"), handlers.UserGetAll)
		users.GET("/:id", middleware.RequirePermission(acc, "user.index", "user.view-all"), handlers.UserGetByID)
		users.PUT("/:id", middleware.RequirePermission(acc, "user.update"), handlers.UserUpdate)
		users.DELETE("/:id", middleware.RequirePermission(acc, "user.delete"), handlers.UserDelete)
	}
}
