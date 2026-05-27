package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
)

func RegisterNotificationRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	notifications := r.Group("/notifications")
	{
		notifications.GET("", handlers.NotificationGetAll)
		notifications.GET("/unread-count", handlers.NotificationCountUnread)
		notifications.GET("/:id", handlers.NotificationGetByID)
		notifications.PATCH("/:id/read", handlers.NotificationMarkAsRead)
		notifications.PATCH("/mark-all-read", handlers.NotificationMarkAllAsRead)
		notifications.DELETE("/:id", handlers.NotificationDelete)
	}
}
