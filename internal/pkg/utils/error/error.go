package error

import (
	"encoding/json"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/spf13/cobra"
)

// APIError represents a structured API error response
type APIError struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
}

// HandleIndexAPIError is a convenience function specifically for index operations
// It extracts the operation from the command context and uses the first argument as index name
func HandleIndexAPIError(err error, cmd *cobra.Command, args []string) {
	if err == nil {
		return
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	// Try to extract JSON error from the error message
	errorMsg := err.Error()

	// Look for JSON-like content in the error message
	var apiErr APIError
	if jsonStart := strings.Index(errorMsg, "{"); jsonStart != -1 {
		jsonContent := errorMsg[jsonStart:]
		if jsonEnd := strings.LastIndex(jsonContent, "}"); jsonEnd != -1 {
			jsonContent = jsonContent[:jsonEnd+1]
			if json.Unmarshal([]byte(jsonContent), &apiErr) == nil && apiErr.Message != "" {
				displayStructuredError(apiErr, verbose)
				return
			}
		}
	}

	// If no structured error found, show the raw error message
	if verbose {
		msg.FailMsg("%s\nFull error: %s\n",
			errorMsg, errorMsg)
	} else {
		msg.FailMsg("%s\n", errorMsg)
	}
}

// displayStructuredError handles structured API error responses
func displayStructuredError(apiErr APIError, verbose bool) {
	// Try to get the message from the body field first (actual API response)
	userMessage := apiErr.Message // fallback to outer message

	// Parse the body field which contains the actual API response
	if apiErr.Body != "" {
		var bodyResponse struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
			Status int `json:"status"`
		}

		if json.Unmarshal([]byte(apiErr.Body), &bodyResponse) == nil && bodyResponse.Error.Message != "" {
			userMessage = bodyResponse.Error.Message
		}
	}

	if userMessage == "" {
		userMessage = "Unknown error occurred"
	}

	if verbose {
		// Show full JSON error in verbose mode - nicely formatted
		jsonBytes, _ := json.MarshalIndent(apiErr, "", "  ")
		msg.FailMsg("%s\n\nFull error response:\n%s\n",
			userMessage, string(jsonBytes))
	} else {
		msg.FailMsg("%s\n", userMessage)
	}
}
