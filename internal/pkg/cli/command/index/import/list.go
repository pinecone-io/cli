package importcmd

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type listImportsCmdOptions struct {
	indexName       string
	limit           int
	paginationToken string
	json            bool
}

// NewListImportsCmd returns the "import list" subcommand.
func NewListImportsCmd() *cobra.Command {
	options := listImportsCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List import operations for an index",
		Long: help.Long(`
			List bulk import operations for the given index, with optional pagination.
		`),
		Example: help.Examples(`
			# List all imports for an index
			pc index import list --index-name my-index

			# List imports with a page size limit
			pc index import list --index-name my-index --limit 10

			# Continue paginating from a previous call
			pc index import list --index-name my-index --pagination-token <token>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)
			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
			if err != nil {
				msg.FailJSON(options.json, "Failed to connect to index: %s\n", err)
				exit.Error(err, "Failed to connect to index")
			}

			err = runListImportsCmd(ctx, ic, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to list imports: %s\n", err)
				exit.Error(err, "Failed to list imports")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "Name of the index to list imports for")
	cmd.Flags().IntVarP(&options.limit, "limit", "l", 0, "Maximum number of imports to return")
	cmd.Flags().StringVarP(&options.paginationToken, "pagination-token", "p", "", "Pagination token to continue a previous listing operation")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runListImportsCmd(ctx context.Context, svc ImportService, options listImportsCmdOptions) error {
	var limit *int32
	if options.limit > 0 {
		l := int32(options.limit)
		limit = &l
	}

	var paginationToken *string
	if options.paginationToken != "" {
		paginationToken = &options.paginationToken
	}

	resp, err := svc.ListImports(ctx, limit, paginationToken)
	if err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(resp))
	} else {
		presenters.PrintImportList(resp)
	}

	return nil
}
