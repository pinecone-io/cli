package index

import (
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/testutils"
)

func TestDeleteCmd_ArgsValidation(t *testing.T) {
	cmd := NewDeleteCmd()

	// Get preset index name validation tests
	tests := testutils.GetIndexNameValidationTests()

	// Use the generic test utility
	testutils.TestCommandArgsAndFlags(t, cmd, tests)
}

func TestDeleteCmd_Usage(t *testing.T) {
	cmd := NewDeleteCmd()

	// Test that the command has proper usage metadata
	testutils.AssertCommandUsage(t, cmd, "delete <name>", "index")
}
