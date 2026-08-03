package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/sashabaranov/go-openai"
)

const generateQueryToolName = "generate_query"

var generateQueryTool = sdk.Tool{
	Type: sdk.ToolTypeFunction,
	Function: &sdk.FunctionDefinition{
		Name:        generateQueryToolName,
		Description: "Generate a single read-only SQL SELECT statement to answer the user's question about attendance/employee data. Only call this when the question needs data from the database.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"sql": {"type": "string", "description": "A single SELECT SQL statement, no trailing semicolon needed."}
			},
			"required": ["sql"]
		}`),
	},
}

// Client is a thin wrapper around the OpenAI chat completion API.
type Client struct {
	api   *sdk.Client
	model string
}

// NewClient creates a Client. apiKey/model/baseURL come from AI_CHAT_API_KEY/
// AI_CHAT_MODEL/AI_CHAT_BASE_URL.
//
// baseURL empty -> OpenAI (api.openai.com). Any OpenAI-wire-compatible
// endpoint works here, e.g. Ollama's local server at http://localhost:11434/v1
// with model "llama3.1" (tool calling requires a tool-capable model).
func NewClient(apiKey, model, baseURL string) *Client {
	cfg := sdk.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	return &Client{api: sdk.NewClientWithConfig(cfg), model: model}
}

// GenerateQuery asks the model to either call generate_query(sql) or answer directly in text.
func (c *Client) GenerateQuery(ctx context.Context, systemPrompt, schemaContext string, history []Message, userMessage string) (*GenerateQueryResult, error) {
	messages := make([]sdk.ChatCompletionMessage, 0, len(history)+2)
	messages = append(messages, sdk.ChatCompletionMessage{
		Role:    sdk.ChatMessageRoleSystem,
		Content: systemPrompt + "\n\n" + schemaContext,
	})
	for _, m := range history {
		role := sdk.ChatMessageRoleUser
		if m.Role == RoleAssistant {
			role = sdk.ChatMessageRoleAssistant
		}
		messages = append(messages, sdk.ChatCompletionMessage{Role: role, Content: m.Content})
	}
	messages = append(messages, sdk.ChatCompletionMessage{Role: sdk.ChatMessageRoleUser, Content: userMessage})

	resp, err := c.api.CreateChatCompletion(ctx, sdk.ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
		Tools:    []sdk.Tool{generateQueryTool},
	})
	if err != nil {
		return nil, fmt.Errorf("openai generate_query call failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("openai returned no choices")
	}

	for _, tc := range resp.Choices[0].Message.ToolCalls {
		if tc.Function.Name != generateQueryToolName {
			continue
		}
		var args struct {
			SQL string `json:"sql"`
		}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			return nil, fmt.Errorf("failed to parse generate_query arguments: %w", err)
		}
		return &GenerateQueryResult{SQL: args.SQL}, nil
	}

	return &GenerateQueryResult{DirectAnswer: resp.Choices[0].Message.Content}, nil
}

// FormatAnswer asks the model to turn query result rows into a natural-language answer.
func (c *Client) FormatAnswer(ctx context.Context, question, rowsJSON string) (string, error) {
	resp, err := c.api.CreateChatCompletion(ctx, sdk.ChatCompletionRequest{
		Model: c.model,
		Messages: []sdk.ChatCompletionMessage{
			{
				Role:    sdk.ChatMessageRoleSystem,
				Content: "Kamu adalah asisten yang menjawab pertanyaan seputar data absensi/karyawan dalam Bahasa Indonesia berdasarkan hasil query. Jawab singkat dan natural, jangan tampilkan SQL atau data mentah.",
			},
			{
				Role:    sdk.ChatMessageRoleUser,
				Content: fmt.Sprintf("Pertanyaan: %s\n\nData hasil query (JSON): %s", question, rowsJSON),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("openai format_answer call failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}
