package backup

import (
	"context"

	backuppresenters "github.com/pinecone-io/cli/internal/pkg/utils/backup/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

type CreateBackupCmdOptions struct {
	json      bool
	name      string
	indexName string
}

func NewCreateBackupCmd() *cobra.Command {
	options := CreateBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a backup from a serverless index",
		Run: func(cmd *cobra.Command, args []string) {
			pc := sdk.NewPineconeClient()
			ctx := context.Background()

			req := &pinecone.CreateBackupParams{
				IndexName: options.indexName,
				Name:      &options.name,
			}
			backup, err := pc.CreateBackup(ctx, req)
			if err != nil {
				msg.FailMsg("Failed to create backup: %s\n", err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(backup)
				pcio.Println(json)
			} else {
				describeCommand := pcio.Sprintf("pc backup describe --id %s", backup.BackupId)
				backupName := "unnamed"
				if backup.Name != nil {
					backupName = *backup.Name
				}
				msg.SuccessMsg("Backup %s created successfully. Run %s to check status. \n\n", style.Emphasis(backupName), style.Code(describeCommand))
				backuppresenters.PrintDescribeBackupTable(backup)
			}
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name you want to give the backup")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&options.indexName, "index", "i", "", "name of the serverless index to backup")
	_ = cmd.MarkFlagRequired("index")

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
