package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
)

func (h *Handlers) AiChatMessage(c *gin.Context) {
	var req dtos.AiChatMessageRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	userID := helpers.GetCallerID(c.Request.Context())
	answer, err := h.svcs.AiChatMessage(c.Request.Context(), userID, req.Message)
	if helpers.HandleError(c, err, "Failed to process AI chat message") {
		return
	}

	helpers.OK(c, "Message processed successfully", dtos.AiChatMessageResponse{Answer: answer})
}

func (h *Handlers) AiChatHistory(c *gin.Context) {
	userID := helpers.GetCallerID(c.Request.Context())
	messages := h.svcs.AiChatHistory(c.Request.Context(), userID)
	helpers.OK(c, "History fetched successfully", dtos.AiChatHistoryResponse{Messages: messages})
}

func (h *Handlers) AiChatReset(c *gin.Context) {
	userID := helpers.GetCallerID(c.Request.Context())
	h.svcs.AiChatReset(c.Request.Context(), userID)
	helpers.OK(c, "Conversation reset successfully", nil)
}
