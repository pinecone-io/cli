package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_ASSISTANT_CHAT_COMPLETIONS = "/assistant/chat/%s/chat/completions"
)

func GetAssistantChatCompletions(asstName string, msg string, stream bool) (*models.ChatCompletionModel, error) {
	outgoingMsg := models.ChatCompletionMessage{
		Role:    "user",
		Content: msg,
	}
	chatHistory := state.ChatHist.Get()
	chat, exists := (*chatHistory.History)[asstName]
	if !exists {
		chat = models.AssistantChat{}
		(*chatHistory.History)[asstName] = chat
	}

	// Add new outgoing messages to existing conversation, this becomes the body
	chat.Messages = append(chat.Messages, outgoingMsg)

	body := models.ChatCompletionRequest{
		Messages: chat.Messages,
		Stream:   stream,
	}

	assistantDataUrl, err := GetAssistantDataBaseUrl()
	if err != nil {
		return nil, err
	}

	var resp *models.ChatCompletionModel
	if !stream {
		resp, err = network.PostAndDecode[models.ChatCompletionRequest, models.ChatCompletionModel](
			assistantDataUrl,
			fmt.Sprintf(URL_ASSISTANT_CHAT_COMPLETIONS, asstName),
			true,
			body,
		)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = network.PostAndStreamChatResponse[models.ChatCompletionRequest](
			assistantDataUrl,
			fmt.Sprintf(URL_ASSISTANT_CHAT_COMPLETIONS, asstName),
			true,
			body,
		)
		if err != nil {
			return nil, err
		}
	}

	// If the request was successful, update the chat history
	chat.Messages = append(chat.Messages, processChatCompletionModel(resp)...)
	(*chatHistory.History)[asstName] = chat
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
