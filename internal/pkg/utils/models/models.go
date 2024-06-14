package models

import "time"

type ChatCompletionRequest struct {
	Stream   bool                    `json:"stream"`
	Messages []ChatCompletionMessage `json:"messages"`
}

type ChatCompletionModel struct {
	Id      string        `json:"id"`
	Choices []ChoiceModel `json:"choices"`
	Model   string        `json:"model"`
}

type ChoiceModel struct {
	FinishReason ChatFinishReason      `json:"finish_reason"`
	Index        int32                 `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AssistantChatHistory map[string]AssistantChat

type AssistantChat struct {
	Messages  []ChatCompletionMessage `json:"messages"`
	CreatedOn time.Time               `json:"created_on"`
}

type ChatFinishReason string

const (
	Stop          ChatFinishReason = "stop"
	Length        ChatFinishReason = "length"
	ContentFilter ChatFinishReason = "content_filter"
	FunctionCall  ChatFinishReason = "function_call"
)

type StreamChatCompletionModel struct {
	Id      string             `json:"id"`
	Choices []ChoiceChunkModel `json:"choices"`
	Model   string             `json:"model"`
}

type StreamChunk struct {
	Data StreamChatCompletionModel `json:"data"`
}

type ChoiceChunkModel struct {
	FinishReason ChatFinishReason      `json:"finish_reason"`
	Index        int32                 `json:"index"`
	Delta        ChatCompletionMessage `json:"delta"`
}

type ContextRefModel struct {
	Id     string   `json:"id"`
	Source string   `json:"source"`
	Text   string   `json:"text"`
	Score  float64  `json:"score"`
	Path   []string `json:"path"`
}
