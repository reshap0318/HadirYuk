package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
)

// RegisterSystemRoutes registers system-level public routes (health, jwks).
func RegisterSystemRoutes(r *gin.Engine, handlers *handlers.Handlers) {
	r.GET("/health", handlers.HealthCheck)
	r.GET("/.well-known/jwks.json", handlers.JWKSGetKeys)
}

// RegisterSystemProtectedRoutes registers system-level protected routes.
func RegisterSystemProtectedRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	r.POST("/upload", handlers.FileUpload)
}
