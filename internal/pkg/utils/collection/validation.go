package collection

import (
	"errors"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

// ValidateCollectionNameArgs validates that exactly one non-empty collection name is provided as a positional argument.
// This is the standard validation used across all collection commands (describe, delete).
func ValidateCollectionNameArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("\b" + style.FailMsg("please provide a collection name"))
	}
	if len(args) > 1 {
		return errors.New("\b" + style.FailMsg("please provide only one collection name"))
	}
	if strings.TrimSpace(args[0]) == "" {
		return errors.New("\b" + style.FailMsg("collection name cannot be empty"))
	}
	return nil
}
