package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

func RegisterLocationRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	locations := r.Group("/locations")
	{
		locations.POST("", middleware.RequirePermission(acc, "location.create"), handlers.LocationCreate)
		locations.GET("", middleware.RequirePermission(acc, "location.index"), handlers.LocationGetAll)
		locations.GET("/:id", middleware.RequirePermission(acc, "location.index"), handlers.LocationGetByID)
		locations.PUT("/:id", middleware.RequirePermission(acc, "location.update"), handlers.LocationUpdate)
		locations.DELETE("/:id", middleware.RequirePermission(acc, "location.delete"), handlers.LocationDelete)
	}
}
