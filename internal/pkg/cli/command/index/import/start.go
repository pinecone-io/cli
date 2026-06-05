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
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type startImportCmdOptions struct {
	indexName     string
	uri           string
	integrationId string
	errorMode     string
	json          bool
}

// NewStartImportCmd returns the "import start" subcommand.
func NewStartImportCmd() *cobra.Command {
	options := startImportCmdOptions{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start an import from a storage URI",
		Long: help.Long(`
			Start an import that loads vector data from a storage provider into a
			serverless index. The URI must begin with the scheme of a supported provider
			(e.g. "s3://").

			For private buckets, configure a storage integration in the Pinecone console
			and pass the integration ID with --integration-id.

			Use --error-mode to control behavior when records fail: "continue" skips
			failing records and keeps going; "abort" stops the import on first error.
			Defaults to "continue".
		`),
		Example: help.Examples(`
			# Start an import from a public S3 bucket
			pc index import start --index-name my-index --uri s3://my-bucket/data/

			# Start an import from a private bucket using a storage integration
			pc index import start --index-name my-index --uri s3://my-bucket/data/ --integration-id intg-123

			# Abort on the first error instead of continuing
			pc index import start --index-name my-index --uri s3://my-bucket/data/ --error-mode abort
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)
			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
			if err != nil {
				msg.FailJSON(options.json, "Failed to connect to index: %s\n", err)
				exit.Error(err, "Failed to connect to index")
			}

			err = runStartImportCmd(ctx, ic, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to start import: %s\n", err)
				exit.Error(err, "Failed to start import")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "Name of the index to import into")
	cmd.Flags().StringVarP(&options.uri, "uri", "u", "", "URI of the data to import (e.g. s3://bucket/path/)")
	cmd.Flags().StringVar(&options.integrationId, "integration-id", "", "Storage integration ID for private buckets")
	cmd.Flags().StringVar(&options.errorMode, "error-mode", "", "How to handle record errors: continue (default) or abort")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("uri")

	return cmd
}

func runStartImportCmd(ctx context.Context, svc ImportService, options startImportCmdOptions) error {
	if strings.TrimSpace(options.uri) == "" {
		return fmt.Errorf("--uri is required")
	}

	var integrationId *string
	if options.integrationId != "" {
		id := options.integrationId
		integrationId = &id
	}

	var errorMode *string
	if options.errorMode != "" {
		em := options.errorMode
		errorMode = &em
	}

	resp, err := svc.StartImport(ctx, options.uri, integrationId, errorMode)
	if err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(resp))
	} else {
		msg.SuccessMsg("Import %s started.\n", style.Emphasis(resp.Id))
		presenters.PrintStartImportTable(resp)
	}

	return nil
}
