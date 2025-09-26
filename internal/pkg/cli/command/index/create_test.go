package index

import (
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/testutils"
)

func TestCreateCmd_ArgsValidation(t *testing.T) {
	cmd := NewCreateIndexCmd()

	// Get preset index name validation tests
	tests := testutils.GetIndexNameValidationTests()

	// Add custom tests for this command (create-specific flags)
	customTests := []testutils.CommandTestConfig{
		{
			Name:        "valid - positional arg with --json flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"json": "true"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --json=false",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"json": "false"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --dimension flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"dimension": "1536"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --metric flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"metric": "cosine"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --cloud and --region flags",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"cloud": "aws", "region": "us-east-1"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --environment flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"environment": "us-east-1-aws"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --pod_type flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"pod_type": "p1.x1"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --model flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"model": "multilingual-e5-large"},
			ExpectError: false,
		},
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
