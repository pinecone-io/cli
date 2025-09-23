package backup

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteBackupCmdOptions struct {
	id string
}

func NewDeleteBackupCmd() *cobra.Command {
	options := DeleteBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a backup",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			err := pc.DeleteBackup(ctx, options.id)
			if err != nil {
				msg.FailMsg("Failed to delete backup %s: %s\n", style.Emphasis(options.id), err)
				exit.Error(err)
			}

			msg.SuccessMsg("Backup %s deleted.\n", style.Emphasis(options.id))
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.id, "id", "i", "", "ID of backup to delete")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
