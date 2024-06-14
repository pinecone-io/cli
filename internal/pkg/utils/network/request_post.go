package network

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func PostAndDecode[B any, R any](baseUrl string, path string, useApiKey bool, body B) (*R, error) {
	return RequestWithBodyAndDecode[B, R](baseUrl, path, http.MethodPost, useApiKey, body)
}

func PostAndStreamChatResponse[B any](baseUrl string, path string, useApiKey bool, body B) (*models.ChatCompletionModel, error) {
	resp, err := RequestWithBody[B](baseUrl, path, http.MethodPost, useApiKey, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var completeResponse string
	var id string
	var model string

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
				fmt.Print(choice.Delta.Content)
				os.Stdout.Sync()
				completeResponse += choice.Delta.Content
			}
			id = chunkResp.Id
			model = chunkResp.Model
		}
	}

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

	if err != nil {
		log.Error().
			Err(err).
			Str("method", http.MethodPost).
			Str("url", baseUrl+path).
			Str("status", resp.Status).
			Msg("Error decoding response")
		return nil, pcio.Errorf("error decoding JSON: %v", err)
	}

	log.Info().
		Str("method", http.MethodPost).
		Str("url", baseUrl+path).
		Msg("Request completed successfully")
	return completionResp, nil
}

func PostMultipartFormDataAndDecode[R any](baseUrl string, path string, useApiKey bool, bodyPath string) (*R, error) {
	url := baseUrl + path

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	file, err := os.Open(bodyPath)
	if err != nil {
		return nil, pcio.Errorf("error opening file: %v", bodyPath)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(bodyPath))
	if err != nil {
		return nil, pcio.Errorf("error creating form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, pcio.Errorf("error copying file to form: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, pcio.Errorf("error closing writer: %v", err)
	}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, pcio.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Info().
		Str("method", "POST").
		Str("url", url).
		Str("multipart/form-data", bodyPath).
		Msg("Sending multipart/form-data request")

	resp, err := performRequest(req, useApiKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "POST").
			Str("url", url)
		return nil, pcio.Errorf("error performing request to %s: %v", url, err)
	}

	var parsedResponse R
	err = decodeResponse(resp, &parsedResponse)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "POST").
			Str("url", url).
			Str("status", resp.Status).
			Msg("Error decoding response")
		return nil, pcio.Errorf("error decoding JSON: %v", err)
	}

	log.Info().
		Str("method", "POST").
		Str("url", url).
		Msg("Request completed successfully")
	return &parsedResponse, nil
}
