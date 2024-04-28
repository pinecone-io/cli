package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func buildRequest(verb string, path string) (*http.Request, error) {
	req, err := http.NewRequest(verb, path, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	if os.Getenv("PINECONE_DEBUG_CURL") == "true" {
		fmt.Printf("curl -X %s %s -H \"Content-Type: application/json\" -H \"User-Agent: Pinecone CLI\" -H \"Authorization: Bearer %s\"\n", verb, path, secrets.OAuth2Token.Get().AccessToken)
	}

	req.Header.Add("User-Agent", "Pinecone CLI")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func performRequest(req *http.Request) (*http.Response, error) {
	// This http client is built using our oauth configurations
	// and is already configured with our access token
	ctx := context.Background()
	client := oauth2.GetHttpClient(ctx)

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "token expired") {
			fmt.Println(err)
			secrets.OAuth2Token.Clear()
			secrets.ConfigFile.Save()
			exit.ErrorMsg(fmt.Sprintf("Your session has expired. Please run %s to log in again.", style.Code("pinecone login")))
		}
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// decodeResponse is a generic function that decodes a JSON HTTP response
// into the provided target type.
func decodeResponse[T any](resp *http.Response, target *T) error {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	return nil
}

func FetchAndDecode[T any](path string, method string) (*T, error) {
	url := DashboardBaseURL + path
	req, err := buildRequest(method, url)
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
		return nil, fmt.Errorf("error building request: %v", err)
	}

	resp, err := performRequest(req)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("url", url)
		return nil, fmt.Errorf("error performing request to %s: %v", url, err)
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
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	log.Info().
		Str("method", method).
		Str("url", url).
		Msg("Request completed successfully")
	return &parsedResponse, nil
}
