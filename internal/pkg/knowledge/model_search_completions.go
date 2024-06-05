package knowledge

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_KNOWLEDGE_MODEL_SEARCH_COMPLETIONS = "/knowledge/chat/%s/chat/completions"
)

type SearchContextModel struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SearchCompletionsRequest struct {
	Messages []SearchContextModel `json:"messages"`
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

func GetKnowledgeModelSearchCompletions(kmName string, content string) (*ChatCompletionModel, error) {
	body := SearchCompletionsRequest{
		Messages: []SearchContextModel{
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	resp, err := network.PostAndDecode[SearchCompletionsRequest, ChatCompletionModel](
		GetKnowledgeDataBaseUrl(),
		fmt.Sprintf(URL_KNOWLEDGE_MODEL_SEARCH_COMPLETIONS, kmName),
		body,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
