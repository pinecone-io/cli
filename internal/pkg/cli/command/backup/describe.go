package backup

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/backup"
	backuppresenters "github.com/pinecone-io/cli/internal/pkg/utils/backup/presenters"
	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeBackupCmdOptions struct {
	json bool
}

func NewDescribeBackupCmd() *cobra.Command {
	options := DescribeBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:          "describe <backup-id>",
		Short:        "Get information on a backup",
		Args:         backup.ValidateBackupIDArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			backupID := args[0]
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			backup, err := pc.DescribeBackup(ctx, backupID)
			if err != nil {
				errorutil.HandleAPIError(err, cmd, args)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(backup)
				fmt.Println(json)
			} else {
				backuppresenters.PrintDescribeBackupTable(backup)
			}
		},
	}

	// No required flags - using positional argument

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
