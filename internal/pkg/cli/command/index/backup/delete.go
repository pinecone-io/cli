package backup

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/confirm"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type deleteBackupCmdOptions struct {
	backupId         string
	skipConfirmation bool
	json             bool
}

func NewDeleteBackupCmd() *cobra.Command {
	options := deleteBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a backup by ID",
		Example: help.Examples(`
			pc index backup delete --id backup-123
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			if !options.skipConfirmation && !options.json {
				confirm.Deletion(
					fmt.Sprintf("This will delete backup %s.", style.Emphasis(options.backupId)),
					"This action cannot be undone.",
				)
			}

			err := runDeleteBackupCmd(ctx, pc, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to delete backup: %s\n", err)
				exit.Error(err, "Failed to delete backup")
			}
		},
	}

	cmd.Flags().StringVarP(&options.backupId, "id", "i", "", "ID of the backup to delete")
	_ = cmd.MarkFlagRequired("id")
	cmd.Flags().BoolVar(&options.skipConfirmation, "skip-confirmation", false, "Skip the deletion confirmation prompt")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output result as JSON (also skips confirmation prompt)")

	return cmd
}

func runDeleteBackupCmd(ctx context.Context, svc BackupService, options deleteBackupCmdOptions) error {
	if strings.TrimSpace(options.backupId) == "" {
		return fmt.Errorf("--id is required")
	}

	if err := svc.DeleteBackup(ctx, options.backupId); err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(struct {
			Deleted bool   `json:"deleted"`
			Id      string `json:"id"`
		}{Deleted: true, Id: options.backupId}))
		return nil
	}

	msg.SuccessMsg("Backup %s deleted.\n", style.Emphasis(options.backupId))
	return nil
}
