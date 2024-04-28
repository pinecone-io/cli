package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
)

func buildRequest(verb string, path string) (*http.Request, error) {
	req, err := http.NewRequest(verb, path, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	if os.Getenv("PINECONE_DEBUG_CURL") == "true" {
		fmt.Printf("curl -X %s %s -H \"Content-Type: application/json\" -H \"User-Agent: Pinecone CLI\" -H \"Authorization: Bearer %s\"\n", verb, path, secrets.AccessToken.Get())
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
		fmt.Println("Error sending request:", err)
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
	if err != nil {
		return nil, fmt.Errorf("error building request: %v", err)
	}

	resp, err := performRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request to %s: %v", url, err)
	}

	var parsedResponse T
	err = decodeResponse(resp, &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &parsedResponse, nil
}
