package namespace

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/spf13/cobra"
)

type deleteNamespaceCmdOptions struct {
	indexName string
	name      string
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
			runDeleteNamespaceCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to delete the namespace from")
	cmd.Flags().StringVar(&options.name, "name", "", "name of the namespace to delete")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runDeleteNamespaceCmd(ctx context.Context, options deleteNamespaceCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	err = ic.DeleteNamespace(ctx, options.name)
	if err != nil {
		msg.FailMsg("Failed to delete namespace: %s", err)
		exit.Error(err, "Failed to delete namespace")
	}
	msg.SuccessMsg("Namespace %s deleted successfully.", options.name)
}
