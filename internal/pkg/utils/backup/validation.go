package backup

import (
	"errors"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

// ValidateBackupIDArgs validates that exactly one non-empty backup ID is provided as a positional argument.
// This is the standard validation used across all backup commands (describe, delete).
func ValidateBackupIDArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("\b" + style.FailMsg("please provide a backup ID"))
	}
	if len(args) > 1 {
		return errors.New("\b" + style.FailMsg("please provide only one backup ID"))
	}
	if strings.TrimSpace(args[0]) == "" {
		return errors.New("\b" + style.FailMsg("backup ID cannot be empty"))
	}
	return nil
}
