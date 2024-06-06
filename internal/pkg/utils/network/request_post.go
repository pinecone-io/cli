package network

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func PostAndDecode[B any, R any](baseUrl string, path string, useApiKey bool, body B) (*R, error) {
	return RequestWithBodyAndDecode[B, R](baseUrl, path, http.MethodPost, useApiKey, body)
}

func PostAndDecodeMultipartFormData[R any](baseUrl string, path string, useApiKey bool, bodyPath string) (*R, error) {
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
