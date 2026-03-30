package namespace

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type deleteNamespaceCmdOptions struct {
	indexName string
	name      string
	json      bool
}

func NewDeleteNamespaceCmd() *cobra.Command {
	options := deleteNamespaceCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a namespace from an index by name",
		Long: help.Long(`
			Delete a namespace from an index.

			Provide the index name and namespace name. 
			WARNING: Deleting a namespace removes its data and cannot be undone.
		`),
		Example: help.Examples(`
			# delete a namespace from an index
			pc index namespace delete --index-name "my-index" --name "tenant-a"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			if strings.TrimSpace(options.indexName) == "" {
				msg.FailJSON(options.json, "Failed to delete namespace: --index-name is required")
				exit.ErrorMsg("Failed to delete namespace: --index-name is required")
			}

			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
			if err != nil {
				msg.FailJSON(options.json, "Failed to delete namespace: %s", err)
				exit.Error(err, "Failed to delete namespace")
			}

			err = runDeleteNamespaceCmd(ctx, ic, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to delete namespace: %s", err)
				exit.Error(err, "Failed to delete namespace")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to delete the namespace from")
	cmd.Flags().StringVar(&options.name, "name", "", "name of the namespace to delete")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output result as JSON")

	return cmd
}

func runDeleteNamespaceCmd(ctx context.Context, ic NamespaceService, options deleteNamespaceCmdOptions) error {
	if strings.TrimSpace(options.name) == "" {
		return fmt.Errorf("--name is required")
	}

	if err := ic.DeleteNamespace(ctx, options.name); err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(struct {
			Deleted   bool   `json:"deleted"`
			Namespace string `json:"namespace"`
			Index     string `json:"index"`
		}{Deleted: true, Namespace: options.name, Index: options.indexName}))
		return nil
	}

	msg.SuccessMsg("Namespace %s deleted successfully.", options.name)
	return nil
}
