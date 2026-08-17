package openai

import "time"

// Role identifies who authored a chat message.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message is a single chat turn, stored in Store and sent back as history/context.
type Message struct {
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// GenerateQueryResult is the outcome of asking the model to answer via generate_query tool call.
// Exactly one of SQL / DirectAnswer is set.
type GenerateQueryResult struct {
	SQL          string
	DirectAnswer string
}
