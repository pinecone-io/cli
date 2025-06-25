package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func buildRequest(verb string, path string, body *bytes.Buffer) (*http.Request, error) {
	if body == nil {
		body = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(verb, path, body)
	if err != nil {
		pcio.Println("Error creating request:", err)
		return nil, err
	}

	if os.Getenv("PINECONE_DEBUG_CURL") == "true" {
		if secrets.OAuth2Token.Get().AccessToken != "" {
			pcio.Printf("curl -X %s %s -H \"Content-Type: application/json\" -H \"User-Agent: Pinecone CLI\" -H \"Authorization: Bearer %s\"\n", verb, path, secrets.OAuth2Token.Get().AccessToken)
		} else {
			pcio.Printf("curl -X %s %s -H \"Content-Type: application/json\" -H \"User-Agent: Pinecone CLI\" -H \"Api-Key: %s\"\n", verb, path, secrets.ApiKey.Get())
		}
	}

	applyHeaders(req, path)
	return req, nil
}

func applyHeaders(req *http.Request, url string) {
	// request-specific headers
	if strings.Contains(url, "assistant") {
		req.Header.Set("X-Project-Id", state.TargetProj.Get().Id)
	}
	if strings.Contains(url, "chat/completions") {
		req.Header.Set("X-Disable-Bearer-Auth", "true")
	}

	// apply to all requests
	req.Header.Add("User-Agent", "Pinecone CLI")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Pinecone-Api-Version", "unstable")
}

func performRequest(req *http.Request) (*http.Response, error) {
	// This http client is built using our oauth configurations
	// and is already configured with our access token
	ctx := context.Background()
	client, err := oauth2.GetHttpClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "token expired") {
			secrets.OAuth2Token.Clear()
			secrets.ConfigFile.Save()
			msg := pcio.Sprintf("Your session has expired. Please run %s to log in again.", style.Code("pinecone login"))
			pcio.Println(msg)
			exit.ErrorMsg(msg)
			return nil, err
		}
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, pcio.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// decodeResponse is a generic function that decodes a JSON HTTP response
// into the provided target type.
func decodeResponse[T any](resp *http.Response, target *T) error {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return pcio.Errorf("error decoding JSON: %v", err)
	}

	return nil
}

func RequestWithBody[B any](baseUrl string, path string, method string, body B) (*http.Response, error) {
	url := baseUrl + path

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("url", url).
			Msg("Error encoding body")
		return nil, pcio.Errorf("error marshalling JSON: %v", err)
	}

	req, err := buildRequest(method, url, &buf)
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
	return resp, nil
}

func RequestWithBodyAndDecode[B any, R any](baseUrl string, path string, method string, body B) (*R, error) {
	url := baseUrl + path

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("url", url).
			Msg("Error encoding body")
		return nil, pcio.Errorf("error marshalling JSON: %v", err)
	}

	req, err := buildRequest(method, url, &buf)
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

func RequestWithoutBodyAndDecode[T any](baseUrl string, path string, method string) (*T, error) {
	url := baseUrl + path

	requestedService := "assistant engine"
	if strings.Contains(url, "console") {
		requestedService = "dashboard"
	}

	req, err := buildRequest(method, url, nil)
	log.Info().
		Str("method", method).
		Str("url", url).
		Msg(fmt.Sprintf("Fetching data from %s", requestedService))
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
