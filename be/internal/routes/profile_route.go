package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
)

// RegisterProfileRoutes registers protected profile routes.
func RegisterProfileRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	me := r.Group("/me")
	{
		me.GET("", handlers.ProfileGet)
		me.PUT("", handlers.ProfileUpdate)
	}
}
