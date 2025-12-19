package backup

import (
	"context"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type describeBackupCmdOptions struct {
	backupId string
	json     bool
}

func NewDescribeBackupCmd() *cobra.Command {
	options := describeBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a backup by ID",
		Example: help.Examples(`
			pc pinecone backup describe --id backup-123
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runDescribeBackupCmd(ctx, pc, options)
			if err != nil {
				msg.FailMsg("Failed to describe backup: %s\n", err)
				exit.Error(err, "Failed to describe backup")
			}
		},
	}

	cmd.Flags().StringVarP(&options.backupId, "id", "i", "", "ID of the backup to describe")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func runDescribeBackupCmd(ctx context.Context, svc BackupService, options describeBackupCmdOptions) error {
	if strings.TrimSpace(options.backupId) == "" {
		return pcio.Errorf("--id is required")
	}

	resp, err := svc.DescribeBackup(ctx, options.backupId)
	if err != nil {
		return err
	}

	if options.json {
		json := text.IndentJSON(resp)
		pcio.Println(json)
	} else {
		presenters.PrintBackupTable(resp)
	}

	return nil
}
