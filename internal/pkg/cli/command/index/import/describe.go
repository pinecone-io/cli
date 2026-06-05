package importcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type describeImportCmdOptions struct {
	indexName string
	importId  string
	json      bool
}

// NewDescribeImportCmd returns the "import describe" subcommand.
func NewDescribeImportCmd() *cobra.Command {
	options := describeImportCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe an import operation by ID",
		Long: help.Long(`
			Show the current status and details of an import operation, including
			percent complete, records imported, and any error messages.
		`),
		Example: help.Examples(`
			pc index import describe --index-name my-index --id import-123
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)
			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
			if err != nil {
				msg.FailJSON(options.json, "Failed to connect to index: %s\n", err)
				exit.Error(err, "Failed to connect to index")
			}

			err = runDescribeImportCmd(ctx, ic, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to describe import: %s\n", err)
				exit.Error(err, "Failed to describe import")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "Name of the index the import belongs to")
	cmd.Flags().StringVarP(&options.importId, "id", "i", "", "ID of the import to describe")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func runDescribeImportCmd(ctx context.Context, svc ImportService, options describeImportCmdOptions) error {
	if strings.TrimSpace(options.importId) == "" {
		return fmt.Errorf("--id is required")
	}

	resp, err := svc.DescribeImport(ctx, options.importId)
	if err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(resp))
	} else {
		presenters.PrintImportTable(resp)
	}

	return nil
}
