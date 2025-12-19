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
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type BackupService interface {
	CreateBackup(ctx context.Context, in *pinecone.CreateBackupParams) (*pinecone.Backup, error)
	DescribeBackup(ctx context.Context, backupId string) (*pinecone.Backup, error)
	ListBackups(ctx context.Context, in *pinecone.ListBackupsParams) (*pinecone.BackupList, error)
	DeleteBackup(ctx context.Context, backupId string) error
	CreateIndexFromBackup(ctx context.Context, in *pinecone.CreateIndexFromBackupParams) (*pinecone.CreateIndexFromBackupResponse, error)
}

type createBackupCmdOptions struct {
	indexName   string
	description string
	name        string
	json        bool
}

func NewCreateBackupCmd() *cobra.Command {
	options := createBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a backup for a serverless index",
		Long: help.Long(`
			Create a backup for a serverless index.

			Provide the index name, and optionally a name/description to identify
			the backup later.
		`),
		Example: help.Examples(`
			pc backup create --index-name my-index
			pc backup create --index-name my-index --name nightly --description "Nightly backup"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runCreateBackupCmd(ctx, pc, options)
			if err != nil {
				msg.FailMsg("Failed to create backup: %s\n", err)
				exit.Error(err, "Failed to create backup")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "i", "", "Name of the index to back up")
	cmd.Flags().StringVarP(&options.description, "description", "d", "", "Optional description for the backup")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Optional name for the backup")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runCreateBackupCmd(ctx context.Context, svc BackupService, options createBackupCmdOptions) error {
	if strings.TrimSpace(options.indexName) == "" {
		return pcio.Errorf("--index-name is required")
	}

	var descPtr *string
	if options.description != "" {
		desc := options.description
		descPtr = &desc
	}
	var namePtr *string
	if options.name != "" {
		n := options.name
		namePtr = &n
	}

	req := &pinecone.CreateBackupParams{
		IndexName:   options.indexName,
		Description: descPtr,
		Name:        namePtr,
	}

	backup, err := svc.CreateBackup(ctx, req)
	if err != nil {
		return err
	}

	if options.json {
		json := text.IndentJSON(backup)
		pcio.Println(json)
	} else {
		msg.SuccessMsg("Backup %s created.\n", styleEmphasisId(backup))
		presenters.PrintBackupTable(backup)
	}

	return nil
}

func styleEmphasisId(backup *pinecone.Backup) string {
	if backup == nil {
		return ""
	}
	return style.Emphasis(backup.BackupId)
}
