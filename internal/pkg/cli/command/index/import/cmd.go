package importcmd

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

// ImportService abstracts the Pinecone IndexConnection for unit testing across import commands.
type ImportService interface {
	StartImport(ctx context.Context, uri string, integrationId *string, errorMode *string) (*pinecone.StartImportResponse, error)
	DescribeImport(ctx context.Context, id string) (*pinecone.Import, error)
	ListImports(ctx context.Context, limit *int32, paginationToken *string) (*pinecone.ListImportsResponse, error)
	CancelImport(ctx context.Context, id string) error
}

var importHelp = help.Long(`
	Manage imports for a serverless index. An import loads vector data from
	a storage provider (S3, GCS, or Azure) directly into an index without requiring you to
	push records through the upsert API. For secure data sources, you can configure a 
	storage integration through the Pinecone console.

	Use these commands to start, describe, list, and cancel import operations.

	Docs:
	  Import data:          https://docs.pinecone.io/guides/index-data/import-data
	  Storage integrations: https://docs.pinecone.io/guides/operations/integrations/manage-storage-integrations
`)

// NewImportCmd returns the parent "import" command with all subcommands attached.
func NewImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "import",
		Short:   "Manage imports for a serverless index",
		Long:    importHelp,
		GroupID: help.GROUP_INDEX_MANAGEMENT.ID,
		Example: help.Examples(`
			# Start an import from an S3 URI
			pc index import start --index-name my-index --uri s3://my-bucket/data/

			# List imports for an index
			pc index import list --index-name my-index

			# Describe a specific import
			pc index import describe --index-name my-index --id import-123

			# Cancel an in-progress import
			pc index import cancel --index-name my-index --id import-123
		`),
	}

	cmd.AddCommand(NewStartImportCmd())
	cmd.AddCommand(NewDescribeImportCmd())
	cmd.AddCommand(NewListImportsCmd())
	cmd.AddCommand(NewCancelImportCmd())

	return cmd
}
