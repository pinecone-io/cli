package backup

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type listBackupsCmdOptions struct {
	indexName       string
	limit           int
	paginationToken string
	json            bool
}

func NewListBackupsCmd() *cobra.Command {
	options := listBackupsCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List backups for the current project",
		Long: help.Long(`
			List backups in the current project, optionally filtered by index name.
		`),
		Example: help.Examples(`
			# List backups for the current project
			pc pinecone backup list

			# List backups for a specific index
			pc pinecone backup list --index-name my-index --limit 10
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runListBackupsCmd(ctx, pc, options)
			if err != nil {
				msg.FailMsg("Failed to list backups: %s\n", err)
				exit.Error(err, "Failed to list backups")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "i", "", "Filter backups by index name")
	cmd.Flags().IntVarP(&options.limit, "limit", "l", 0, "Maximum number of backups to return")
	cmd.Flags().StringVarP(&options.paginationToken, "pagination-token", "p", "", "Pagination token to continue a previous listing operation")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runListBackupsCmd(ctx context.Context, svc BackupService, options listBackupsCmdOptions) error {
	var limit *int
	if options.limit > 0 {
		limit = &options.limit
	}

	var paginationToken *string
	if options.paginationToken != "" {
		paginationToken = &options.paginationToken
	}

	var indexName *string
	if options.indexName != "" {
		indexName = &options.indexName
	}

	resp, err := svc.ListBackups(ctx, &pinecone.ListBackupsParams{
		IndexName:       indexName,
		Limit:           limit,
		PaginationToken: paginationToken,
	})
	if err != nil {
		return err
	}

	if options.json {
		json := text.IndentJSON(resp)
		pcio.Println(json)
	} else {
		presenters.PrintBackupList(resp)
	}

	return nil
}
