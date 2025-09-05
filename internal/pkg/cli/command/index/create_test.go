package index

import (
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/testutils"
)

func TestCreateCmd_ArgsValidation(t *testing.T) {
	cmd := NewCreateIndexCmd()

	// Get preset index name validation tests
	tests := testutils.GetIndexNameValidationTests()

	// Add custom tests for this command (create-specific business logic)
	customTests := []testutils.CommandTestConfig{
		{
			Name:        "error - no arguments but with --json flag",
			Args:        []string{},
			Flags:       map[string]string{"json": "true"},
			ExpectError: true,
			ErrorSubstr: "please provide an index name",
		},
		{
			Name:        "error - multiple arguments with --json flag",
			Args:        []string{"index1", "index2"},
			Flags:       map[string]string{"json": "true"},
			ExpectError: true,
			ErrorSubstr: "please provide only one index name",
		},
	}

	// Combine preset tests with custom tests
	allTests := append(tests, customTests...)

	// Use the generic test utility
	testutils.TestCommandArgsAndFlags(t, cmd, allTests)
}

func TestCreateCmd_Flags(t *testing.T) {
	cmd := NewCreateIndexCmd()

	// Test that the command has a properly configured --json flag
	testutils.AssertJSONFlag(t, cmd)
}

func TestCreateCmd_Usage(t *testing.T) {
	cmd := NewCreateIndexCmd()

	// Test that the command has proper usage metadata
	testutils.AssertCommandUsage(t, cmd, "create <name>", "index")
}
