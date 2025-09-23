package backup

import (
	"context"

	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

type RestoreBackupCmdOptions struct {
	json      bool
	indexName string
	backupID  string
}

func NewRestoreBackupCmd() *cobra.Command {
	options := RestoreBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Create a new index from a backup",
		Run: func(cmd *cobra.Command, args []string) {
			pc := sdk.NewPineconeClient()
			ctx := context.Background()

			req := &pinecone.CreateIndexFromBackupParams{
				Name:     options.indexName,
				BackupId: options.backupID,
			}
			restoreJob, err := pc.CreateIndexFromBackup(ctx, req)
			if err != nil {
				errorutil.HandleAPIError(err, cmd, args)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(restoreJob)
				pcio.Println(json)
			} else {
				describeCommand := pcio.Sprintf("pc index describe --name %s", options.indexName)
				msg.SuccessMsg("Index %s restore job initiated successfully. Run %s to check status.\n", style.Emphasis(options.indexName), style.Code(describeCommand))
			}
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.indexName, "index", "i", "", "name for the new index to create")
	_ = cmd.MarkFlagRequired("index")
	cmd.Flags().StringVarP(&options.backupID, "backup-id", "b", "", "ID of the backup to restore from")
	_ = cmd.MarkFlagRequired("backup-id")

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
