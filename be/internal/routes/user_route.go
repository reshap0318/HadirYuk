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
		users.GET("", middleware.RequirePermission(acc, "user.index"), handlers.UserGetAll)
		users.GET("/:id", middleware.RequirePermission(acc, "user.index"), handlers.UserGetByID)
		users.PUT("/:id", middleware.RequirePermission(acc, "user.edit"), handlers.UserUpdate)
		users.DELETE("/:id", middleware.RequirePermission(acc, "user.delete"), handlers.UserDelete)
	}

	me := r.Group("/me")
	{
		me.GET("", handlers.ProfileGet)
		me.PUT("", handlers.ProfileUpdate)
	}
}
