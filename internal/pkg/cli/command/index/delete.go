package index

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/confirm"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type deleteCmdOptions struct {
	indexName        string
	skipConfirmation bool
	json             bool
}

// DeleteIndexService abstracts the Pinecone Go SDK for unit testing (runDeleteIndexCmd)
type DeleteIndexService interface {
	DeleteIndex(ctx context.Context, name string) error
}

func NewDeleteCmd() *cobra.Command {
	options := deleteCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an index by name",
		Example: help.Examples(`
			pc index delete --index-name "index-name"
		`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("index-name") && !cmd.Flags().Changed("name") {
				return fmt.Errorf("required flag(s) \"index-name\" not set")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			if !options.skipConfirmation && !options.json {
				confirm.Deletion(
					fmt.Sprintf("This will delete index %s and all of its data.", style.Emphasis(options.indexName)),
					"This action cannot be undone.",
				)
			}

			err := runDeleteIndexCmd(ctx, pc, options)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailJSON(options.json, "The index %s does not exist\n", style.Emphasis(options.indexName))
					exit.Errorf(err, "The index %s does not exist", style.Emphasis(options.indexName))
				} else {
					msg.FailJSON(options.json, "Failed to delete index %s: %s\n", style.Emphasis(options.indexName), err)
					exit.Errorf(err, "Failed to delete index %s", style.Emphasis(options.indexName))
				}
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.indexName, "index-name", "i", "", "name of index to delete")
	cmd.Flags().StringVarP(&options.indexName, "name", "n", "", "name of index to delete")
	_ = cmd.Flags().MarkDeprecated("name", "use --index-name instead")
	cmd.Flags().BoolVar(&options.skipConfirmation, "skip-confirmation", false, "Skip the deletion confirmation prompt")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output result as JSON (also skips confirmation prompt)")

	return cmd
}

func runDeleteIndexCmd(ctx context.Context, svc DeleteIndexService, options deleteCmdOptions) error {
	if err := svc.DeleteIndex(ctx, options.indexName); err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(struct {
			Deleted bool   `json:"deleted"`
			Name    string `json:"name"`
		}{Deleted: true, Name: options.indexName}))
		return nil
	}

	msg.SuccessMsg("Index %s deleted.\n", style.Emphasis(options.indexName))
	return nil
}
