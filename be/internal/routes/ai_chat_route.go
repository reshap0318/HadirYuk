package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

func RegisterAiChatRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	aiChat := r.Group("/ai-chat")
	{
		aiChat.POST("/message", middleware.RequirePermission(acc, "ai-chat.query"), handlers.AiChatMessage)
		aiChat.GET("/history", middleware.RequirePermission(acc, "ai-chat.query"), handlers.AiChatHistory)
		aiChat.POST("/reset", middleware.RequirePermission(acc, "ai-chat.query"), handlers.AiChatReset)
	}
}
