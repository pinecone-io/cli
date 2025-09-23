package backup

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/backup"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewDeleteBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "delete <backup-id>",
		Short:        "Delete a backup",
		Args:         backup.ValidateBackupIDArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			backupID := args[0]
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			err := pc.DeleteBackup(ctx, backupID)
			if err != nil {
				msg.FailMsg("Failed to delete backup %s: %s\n", style.Emphasis(backupID), err)
				exit.Error(err)
			}

			msg.SuccessMsg("Backup %s deleted.\n", style.Emphasis(backupID))
		},
	}

	// No flags needed - using positional argument

	return cmd
}
