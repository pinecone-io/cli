package dashboard

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func FetchAndDecode[T any](path string) (*T, error) {
	url := DashboardBaseURL + path
	method := http.MethodGet
	req, err := buildRequest(method, url, nil)
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

	var parsedResponse T
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
