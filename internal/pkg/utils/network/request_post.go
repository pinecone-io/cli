package network

import (
	"encoding/json"
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func PostAndDecode[B any, R any](baseUrl string, path string, body B) (*R, error) {
	return RequestWithBodyAndDecode[B, R](baseUrl, path, http.MethodPost, body)
}

func RequestWithBodyAndDecode[B any, R any](baseUrl string, path string, method string, body B) (*R, error) {
	url := baseUrl + path

	var bodyJson []byte
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("url", url).
			Msg("Error marshalling JSON")
		return nil, pcio.Errorf("error marshalling JSON: %v", err)
	}

	req, err := buildRequest(method, url, bodyJson)
	log.Info().
		Str("method", method).
		Str("url", url).
		Msg("Fetching data from dashboard")
	if err != nil {
		log.Error().
			Err(err).
			Str("url", url).
			Str("method", method).
			Msg("Error building request")
		return nil, pcio.Errorf("error building request: %v", err)
	}

	resp, err := performRequest(req)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("url", url)
		return nil, pcio.Errorf("error performing request to %s: %v", url, err)
	}

	var parsedResponse R
	err = decodeResponse(resp, &parsedResponse)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("url", url).
			Str("status", resp.Status).
			Msg("Error decoding response")
		return nil, pcio.Errorf("error decoding JSON: %v", err)
	}

	log.Info().
		Str("method", method).
		Str("url", url).
		Msg("Request completed successfully")
	return &parsedResponse, nil
}
