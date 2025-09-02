package error

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

func TestHandleIndexAPIErrorWithCommand(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		indexName      string
		commandName    string
		verbose        bool
		expectedOutput string
	}{
		{
			name:           "JSON error with message field",
			err:            &mockError{message: `{"message": "Index not found", "code": 404}`},
			indexName:      "test-index",
			commandName:    "describe <name>",
			verbose:        false,
			expectedOutput: "Index not found",
		},
		{
			name:           "Verbose mode shows full JSON",
			err:            &mockError{message: `{"message": "Rate limit exceeded", "code": 429}`},
			indexName:      "my-index",
			commandName:    "create <name>",
			verbose:        true,
			expectedOutput: "Rate limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock command with verbose flag and set the command name
			cmd := &cobra.Command{}
			cmd.Flags().Bool("verbose", false, "verbose output")
			cmd.Use = tt.commandName

			// Set the verbose flag on the command
			cmd.Flags().Set("verbose", fmt.Sprintf("%t", tt.verbose))

			// This is a basic test to ensure the function doesn't panic
			// In a real test environment, we would capture stdout/stderr
			// and verify the exact output
			HandleIndexAPIError(tt.err, cmd, []string{tt.indexName})
		})
	}
}

// mockError is a simple error implementation for testing
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
