package index

import (
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/testutils"
)

func TestConfigureCmd_ArgsValidation(t *testing.T) {
	cmd := NewConfigureIndexCmd()

	// Get preset index name validation tests
	tests := testutils.GetIndexNameValidationTests()

	// Add custom tests for this command (configure-specific flags)
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
			Name:        "valid - positional arg with --pod_type flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"pod_type": "p1.x1"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --replicas flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"replicas": "2"},
			ExpectError: false,
		},
		{
			Name:        "valid - positional arg with --deletion_protection flag",
			Args:        []string{"my-index"},
			Flags:       map[string]string{"deletion_protection": "enabled"},
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

func TestConfigureCmd_Flags(t *testing.T) {
	cmd := NewConfigureIndexCmd()

	// Test that the command has a properly configured --json flag
	testutils.AssertJSONFlag(t, cmd)
}

func TestConfigureCmd_Usage(t *testing.T) {
	cmd := NewConfigureIndexCmd()

	// Test that the command has proper usage metadata
	testutils.AssertCommandUsage(t, cmd, "configure <name>", "index")
}
