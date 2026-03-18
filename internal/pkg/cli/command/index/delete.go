package index

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

type deleteCmdOptions struct {
	name string
	json bool
}

type deleteIndexService interface {
	DeleteIndex(ctx context.Context, name string) error
}

func NewDeleteCmd() *cobra.Command {
	options := deleteCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an index by name",
		Example: help.Examples(`
			pc index delete --name "index-name"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runDeleteIndexCmd(ctx, pc, options)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailMsg("The index %s does not exist\n", style.Emphasis(options.name))
					exit.Errorf(err, "The index %s does not exist", style.Emphasis(options.name))
				} else {
					msg.FailMsg("Failed to delete index %s: %s\n", style.Emphasis(options.name), err)
					exit.Errorf(err, "Failed to delete index %s", style.Emphasis(options.name))
				}
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to delete")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output result as JSON")

	return cmd
}

func runDeleteIndexCmd(ctx context.Context, svc deleteIndexService, options deleteCmdOptions) error {
	if err := svc.DeleteIndex(ctx, options.name); err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(struct {
			Deleted bool   `json:"deleted"`
			Name    string `json:"name"`
		}{Deleted: true, Name: options.name}))
		return nil
	}

	msg.SuccessMsg("Index %s deleted.\n", style.Emphasis(options.name))
	return nil
}
