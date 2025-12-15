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
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type createNamespaceCmdOptions struct {
	indexName      string
	name           string
	metadataSchema []string
	json           bool
}

func NewCreateNamespaceCmd() *cobra.Command {
	options := createNamespaceCmdOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new namespace in an index",
		Long:    help.Long(``),
		Example: help.Examples(``),
		Run: func(cmd *cobra.Command, args []string) {
			runCreateNamespaceCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to create the namespace in")
	cmd.Flags().StringVar(&options.name, "name", "", "name of the namespace to create")
	cmd.Flags().StringSliceVar(&options.metadataSchema, "schema", []string{}, "metadata schema for the namespace")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runCreateNamespaceCmd(ctx context.Context, options createNamespaceCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	req := &pinecone.CreateNamespaceParams{
		Name:   options.name,
		Schema: sdk.BuildMetadataSchema(options.metadataSchema),
	}
	ns, err := ic.CreateNamespace(ctx, req)
	if err != nil {
		msg.FailMsg("Failed to create namespace: %s", err)
		exit.Error(err, "Failed to create namespace")
	}

	if options.json {
		json := text.IndentJSON(ns)
		pcio.Println(json)
	} else {
		msg.SuccessMsg("Namespace %s created successfully.", options.name)
		presenters.PrintDescribeNamespaceTable(ns)
	}
}
