package dtos

import "time"

// AiChatMessageRequest is the payload for POST /api/ai-chat/message.
type AiChatMessageRequest struct {
	Message string `json:"message" validate:"required,max=500"`
}

// AiChatMessageResponse wraps the AI's natural-language answer.
type AiChatMessageResponse struct {
	Answer string `json:"answer"`
}

// AiChatMessageDTO is one chat turn, used for scrollback in the history endpoint.
type AiChatMessageDTO struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// AiChatHistoryResponse wraps GET /api/ai-chat/history.
type AiChatHistoryResponse struct {
	Messages []AiChatMessageDTO `json:"messages"`
}
