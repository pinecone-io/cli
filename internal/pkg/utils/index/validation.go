package index

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

// ValidateIndexNameArgs validates that exactly one non-empty index name is provided as a positional argument.
// This is the standard validation used across all index commands (create, describe, delete, configure).
func ValidateIndexNameArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("please provide an index name")
	}
	if len(args) > 1 {
		return errors.New("please provide only one index name")
	}
	if strings.TrimSpace(args[0]) == "" {
		return errors.New("index name cannot be empty")
	}
	return nil
}
