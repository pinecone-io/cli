package backup

import (
	"context"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type deleteBackupCmdOptions struct {
	backupId string
}

func NewDeleteBackupCmd() *cobra.Command {
	options := deleteBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a backup by ID",
		Example: help.Examples(`
			pc pinecone backup delete --id backup-123
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runDeleteBackupCmd(ctx, pc, options)
			if err != nil {
				msg.FailMsg("Failed to delete backup: %s\n", err)
				exit.Error(err, "Failed to delete backup")
			}
		},
	}

	cmd.Flags().StringVarP(&options.backupId, "id", "i", "", "ID of the backup to delete")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func runDeleteBackupCmd(ctx context.Context, svc BackupService, options deleteBackupCmdOptions) error {
	if strings.TrimSpace(options.backupId) == "" {
		return pcio.Errorf("--id is required")
	}

	if err := svc.DeleteBackup(ctx, options.backupId); err != nil {
		return err
	}

	msg.SuccessMsg("Backup %s deleted.\n", style.Emphasis(options.backupId))
	return nil
}
