package importcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type cancelImportCmdOptions struct {
	indexName string
	importId  string
	json      bool
}

// NewCancelImportCmd returns the "import cancel" subcommand.
func NewCancelImportCmd() *cobra.Command {
	options := cancelImportCmdOptions{}

	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel an in-progress import operation",
		Long: help.Long(`
			Cancel an import operation that is currently pending or in progress.
			Already-completed or failed imports cannot be cancelled.
		`),
		Example: help.Examples(`
			pc index import cancel --index-name my-index --id import-123
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)
			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
			if err != nil {
				msg.FailJSON(options.json, "Failed to connect to index: %s\n", err)
				exit.Error(err, "Failed to connect to index")
			}

			err = runCancelImportCmd(ctx, ic, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to cancel import: %s\n", err)
				exit.Error(err, "Failed to cancel import")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "Name of the index the import belongs to")
	cmd.Flags().StringVarP(&options.importId, "id", "i", "", "ID of the import to cancel")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func runCancelImportCmd(ctx context.Context, svc ImportService, options cancelImportCmdOptions) error {
	if strings.TrimSpace(options.importId) == "" {
		return fmt.Errorf("--id is required")
	}

	if err := svc.CancelImport(ctx, options.importId); err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(struct {
			Cancelled bool   `json:"cancelled"`
			Id        string `json:"id"`
		}{Cancelled: true, Id: options.importId}))
		return nil
	}

	msg.SuccessMsg("Import %s cancelled.\n", style.Emphasis(options.importId))
	return nil
}
