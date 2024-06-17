package assistants

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_ASSISTANT_CHAT_COMPLETIONS = "/assistant/chat/%s/chat/completions"
)

const maxLineWidth = 80

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

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	var resp *models.ChatCompletionModel
	if !stream {
		resp, err = network.PostAndDecode[models.ChatCompletionRequest, models.ChatCompletionModel](
			assistantDataUrl,
			fmt.Sprintf(URL_ASSISTANT_CHAT_COMPLETIONS, asstName),
			body,
		)
		s.Stop()
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = PostAndStreamChatResponse[models.ChatCompletionRequest](
			assistantDataUrl,
			fmt.Sprintf(URL_ASSISTANT_CHAT_COMPLETIONS, asstName),
			body,
			s,
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

func PostAndStreamChatResponse[B any](baseUrl string, path string, body B, spinner *spinner.Spinner) (*models.ChatCompletionModel, error) {
	resp, err := network.RequestWithBody[B](baseUrl, path, http.MethodPost, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	spinner.Stop()

	var completeResponse string
	var id string
	var model string
	var currentLineLength int

	// stream response and print as we go
	fmt.Print("\n")
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			dataStr := strings.TrimPrefix(line, "data:")
			dataStr = strings.TrimSpace(dataStr)

			var chunkResp *models.StreamChatCompletionModel
			if err := json.Unmarshal([]byte(dataStr), &chunkResp); err != nil {
				return nil, pcio.Errorf("error unmarshaling chunk: %v", err)
			}

			for _, choice := range chunkResp.Choices {
				content := choice.Delta.Content
				contentLength := len(content)

				// would exceed max line length, push to new line
				if currentLineLength+contentLength > maxLineWidth {
					fmt.Print("\n")
					currentLineLength = 0
				}
				fmt.Print(choice.Delta.Content)
				os.Stdout.Sync()
				currentLineLength += contentLength
				completeResponse += choice.Delta.Content
			}
			id = chunkResp.Id
			model = chunkResp.Model
		}
	}
	fmt.Print("\n")

	completionResp := &models.ChatCompletionModel{
		Id:    id,
		Model: model,
		Choices: []models.ChoiceModel{
			{
				FinishReason: "stop",
				Index:        0,
				Message: models.ChatCompletionMessage{
					Content: completeResponse,
					Role:    "assistant",
				},
			},
		},
	}

	log.Info().
		Str("method", http.MethodPost).
		Str("url", baseUrl+path).
		Msg("Request completed successfully")
	return completionResp, nil
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
