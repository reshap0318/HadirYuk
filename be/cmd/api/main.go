package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/reshap0318/hadirYuk/internal/di"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
	"github.com/reshap0318/hadirYuk/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	host := helpers.GetEnv("APP_HOST", "0.0.0.0")
	port := helpers.GetEnv("APP_PORT", "8080")
	trustedProxies := helpers.GetEnv("TRUSTED_PROXIES", "")
	allowedOrigins := helpers.GetEnv("ALLOWED_ORIGINS", "*")

	gin.SetMode(helpers.GetEnv("GIN_MODE", "release"))

	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	r := gin.Default()

	if trustedProxies != "" {
		if err := r.SetTrustedProxies(strings.Split(trustedProxies, ",")); err != nil {
			log.Printf("Warning: failed to set trusted proxies: %v", err)
		}
	}

	r.Use(middleware.RateLimit(container.RateLimiter))
	r.Use(middleware.CORS(allowedOrigins))

	r.Static("/storage", "./storage")

	apiGroup := r.Group("/api")
	{
		routes.RegisterSystemRoutes(r, container.Handlers)
		routes.RegisterAuthRoutes(apiGroup, container.Handlers)
	}

	protected := apiGroup.Group("")
	protected.Use(middleware.JWTAuth(container.Services))
	{
		routes.RegisterAuthProtectedRoutes(protected, container.Handlers)
		routes.RegisterPermissionRoutes(protected, container.Handlers, container.Access)
		routes.RegisterRoleRoutes(protected, container.Handlers, container.Access)
		routes.RegisterUserRoutes(protected, container.Handlers, container.Access)
		routes.RegisterProfileRoutes(protected, container.Handlers, container.Access)
		routes.RegisterNotificationRoutes(protected, container.Handlers)
		routes.RegisterSystemProtectedRoutes(protected, container.Handlers)
		routes.RegisterShiftRoutes(protected, container.Handlers, container.Access)
		routes.RegisterShiftAssignmentRoutes(protected, container.Handlers, container.Access)
		routes.RegisterLocationRoutes(protected, container.Handlers, container.Access)
		routes.RegisterLeaveTypeRoutes(protected, container.Handlers, container.Access)
		routes.RegisterFaceRoutes(protected, container.Handlers, container.Access)
	}

	addr := host + ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
