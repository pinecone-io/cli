package models

import "time"

type ChatCompletionRequest struct {
	Messages []ChatCompletionMessage `json:"messages"`
}

type ChatCompletionModel struct {
	Id      string        `json:"id"`
	Choices []ChoiceModel `json:"choices"`
	Model   string        `json:"model"`
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

type ChoiceModel struct {
	FinishReason ChatFinishReason      `json:"finish_reason"`
	Index        int32                 `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
}
