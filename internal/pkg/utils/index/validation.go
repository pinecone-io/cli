package index

import (
	"errors"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

// ValidateIndexNameArgs validates that exactly one non-empty index name is provided as a positional argument.
// This is the standard validation used across all index commands (create, describe, delete, configure).
func ValidateIndexNameArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("\b" + style.FailMsg("please provide an index name"))
	}
	if len(args) > 1 {
		return errors.New("\b" + style.FailMsg("please provide only one index name"))
	}
	if strings.TrimSpace(args[0]) == "" {
		return errors.New("\b" + style.FailMsg("index name cannot be empty"))
	}
	return nil
}
