package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

// RegisterFaceRoutes registers face photo routes.
func RegisterFaceRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	// Admin face photo management under /api/users/:id
	users := r.Group("/users")
	{
		users.PUT("/:id/face-photo", middleware.RequirePermission(acc, "user.update"), handlers.AdminFacePhotoUpload)
		users.DELETE("/:id/face-photo", middleware.RequirePermission(acc, "user.update"), handlers.AdminFacePhotoDelete)
	}

	// Face recognition for attendance
	face := r.Group("/face")
	{
		face.POST("/match", handlers.FaceMatch)
	}
}
