package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_ASSISTANT_CHAT_COMPLETIONS         = "/knowledge/chat/%s/chat/completions"
	URL_ASSISTANT_CHAT_COMPLETIONS_STAGING = "/assistant/chat/%s/chat/completions"
)

func getAssistantChatCompletionsUrl() string {
	if config.Environment.Get() == "production" {
		return URL_ASSISTANT_CHAT_COMPLETIONS
	} else {
		return URL_ASSISTANT_CHAT_COMPLETIONS_STAGING
	}
}

func GetAssistantChatCompletions(kmName string, msg string) (*models.ChatCompletionModel, error) {
	outgoingMsg := models.ChatCompletionMessage{
		Role:    "user",
		Content: msg,
	}
	chatHistory := state.ChatHist.Get()
	chat, exists := (*chatHistory.History)[kmName]
	if !exists {
		chat = models.AssistantChat{}
		(*chatHistory.History)[kmName] = chat
	}

	// Add new outgoing messages to existing conversation, this becomes the body
	chat.Messages = append(chat.Messages, outgoingMsg)

	body := models.ChatCompletionRequest{
		Messages: chat.Messages,
	}

	assistantDataUrl, err := GetAssistantDataBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.PostAndDecode[models.ChatCompletionRequest, models.ChatCompletionModel](
		assistantDataUrl,
		fmt.Sprintf(getAssistantChatCompletionsUrl(), kmName),
		true,
		body,
	)
	if err != nil {
		return nil, err
	}

	// If the request was successful, update the chat history
	chat.Messages = append(chat.Messages, processChatCompletionModel(resp)...)
	(*chatHistory.History)[kmName] = chat
	state.ChatHist.Set(&chatHistory)

	return resp, nil
}

func processChatCompletionModel(resp *models.ChatCompletionModel) []models.ChatCompletionMessage {
	var messages []models.ChatCompletionMessage

	log.Trace().
		Str("Id", resp.Id).
		Str("Model", resp.Model).
		Msg("processing ChatCompletionModel")

	for _, choice := range resp.Choices {
		log.Trace().
			Str("Message", choice.Message.Content).
			Int32("Index", choice.Index).
			Str("FinishReason", string(choice.FinishReason)).
			Msg("found ChoiceModel")

		messages = append(messages, choice.Message)
	}

	return messages
}
