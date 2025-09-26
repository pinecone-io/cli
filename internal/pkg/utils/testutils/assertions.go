package testutils

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// AssertCommandUsage tests that a command has proper usage metadata
// This function can be reused by any command to test its usage information
func AssertCommandUsage(t *testing.T, cmd *cobra.Command, expectedUsage string, expectedDomain string) {
	t.Helper()

	// Test usage string
	if cmd.Use != expectedUsage {
		t.Errorf("expected Use to be %q, got %q", expectedUsage, cmd.Use)
	}

	// Test short description exists
	if cmd.Short == "" {
		t.Error("expected command to have a short description")
	}

	// Test description mentions domain
	if !strings.Contains(strings.ToLower(cmd.Short), expectedDomain) {
		t.Errorf("expected short description to mention %q, got %q", expectedDomain, cmd.Short)
	}
}

// AssertJSONFlag tests the common --json flag pattern used across commands
// This function comprehensively tests flag definition, properties, and functionality
// This function can be reused by any command that has a --json flag
func AssertJSONFlag(t *testing.T, cmd *cobra.Command) {
	t.Helper()

	// Test that the json flag is properly defined
	jsonFlag := cmd.Flags().Lookup("json")
	if jsonFlag == nil {
		t.Error("expected --json flag to be defined")
		return
	}

	// Test that it's a boolean flag
	if jsonFlag.Value.Type() != "bool" {
		t.Errorf("expected --json flag to be bool type, got %s", jsonFlag.Value.Type())
	}

	// Test that the flag is optional (not required)
	if jsonFlag.Annotations[cobra.BashCompOneRequiredFlag] != nil {
		t.Error("expected --json flag to be optional")
	}

	// Test that the flag has a description
	if jsonFlag.Usage == "" {
		t.Error("expected --json flag to have a usage description")
	}

	// Test that the description mentions JSON
	if !strings.Contains(strings.ToLower(jsonFlag.Usage), "json") {
		t.Errorf("expected --json flag description to mention 'json', got %q", jsonFlag.Usage)
	}

	// Test setting json flag to true
	err := cmd.Flags().Set("json", "true")
	if err != nil {
		t.Errorf("failed to set --json flag to true: %v", err)
	}

	jsonValue, err := cmd.Flags().GetBool("json")
	if err != nil {
		t.Errorf("failed to get --json flag value: %v", err)
	}
	if !jsonValue {
		t.Error("expected --json flag to be true after setting it")
	}

	// Test setting json flag to false
	err = cmd.Flags().Set("json", "false")
	if err != nil {
		t.Errorf("failed to set --json flag to false: %v", err)
	}

	jsonValue, err = cmd.Flags().GetBool("json")
	if err != nil {
		t.Errorf("failed to get --json flag value: %v", err)
	}
	if jsonValue {
		t.Error("expected --json flag to be false after setting it")
	}
}
