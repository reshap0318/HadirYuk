package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
)

// RegisterAuthRoutes registers public authentication routes.
func RegisterAuthRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", handlers.AuthLogin)
		auth.POST("/refresh", handlers.AuthRefreshToken)
		auth.POST("/forgot-password", handlers.AuthForgetPassword)
		auth.POST("/reset-password", handlers.AuthResetPassword)
	}
}

// RegisterAuthProtectedRoutes registers protected authentication routes.
func RegisterAuthProtectedRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/logout", handlers.AuthLogout)
	}
}
