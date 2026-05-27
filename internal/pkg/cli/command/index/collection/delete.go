package collection

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type deleteCollectionCmdOptions struct {
	name string
	json bool
}

// DeleteCollectionService abstracts the Pinecone Go SDK for unit testing (runDeleteCollectionCmd)
type DeleteCollectionService interface {
	DeleteCollection(ctx context.Context, name string) error
}

func NewDeleteCollectionCmd() *cobra.Command {
	options := deleteCollectionCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a collection",
		Example: help.Examples(`
			pc index collection delete --name my-collection
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runDeleteCollectionCmd(ctx, pc, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to delete collection %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err, "Failed to delete collection")
			}
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of collection to delete")
	_ = cmd.MarkFlagRequired("name")

	// Optional flags
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output result as JSON")

	return cmd
}

func runDeleteCollectionCmd(ctx context.Context, svc DeleteCollectionService, options deleteCollectionCmdOptions) error {
	if err := svc.DeleteCollection(ctx, options.name); err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(struct {
			Deleted bool   `json:"deleted"`
			Name    string `json:"name"`
		}{Deleted: true, Name: options.name}))
		return nil
	}

	msg.SuccessMsg("Collection %s deleted.\n", style.Emphasis(options.name))
	return nil
}
