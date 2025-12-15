package namespace

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type describeNamespaceCmdOptions struct {
	indexName string
	name      string
	json      bool
}

func NewDescribeNamespaceCmd() *cobra.Command {
	options := describeNamespaceCmdOptions{}

	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Describe a namespace from an index by name",
		Long:    help.Long(``),
		Example: help.Examples(``),
		Run: func(cmd *cobra.Command, args []string) {
			runDescribeNamespaceCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to describe the namespace from")
	cmd.Flags().StringVar(&options.name, "name", "", "name of the namespace to describe")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runDescribeNamespaceCmd(ctx context.Context, options describeNamespaceCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	ns, err := ic.DescribeNamespace(ctx, options.name)
	if err != nil {
		msg.FailMsg("Failed to describe namespace: %s", err)
		exit.Error(err, "Failed to describe namespace")
	}

	if options.json {
		json := text.IndentJSON(ns)
		pcio.Println(json)
	} else {
		presenters.PrintDescribeNamespaceTable(ns)
	}
}
