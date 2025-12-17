package namespace

import (
	"context"
	"strings"

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

type NamespaceService interface {
	CreateNamespace(ctx context.Context, req *pinecone.CreateNamespaceParams) (*pinecone.NamespaceDescription, error)
	DescribeNamespace(ctx context.Context, name string) (*pinecone.NamespaceDescription, error)
	ListNamespaces(ctx context.Context, params *pinecone.ListNamespacesParams) (*pinecone.ListNamespacesResponse, error)
	DeleteNamespace(ctx context.Context, name string) error
}

type createNamespaceCmdOptions struct {
	indexName      string
	name           string
	metadataSchema []string
	json           bool
}

func NewCreateNamespaceCmd() *cobra.Command {
	options := createNamespaceCmdOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace in an index",
		Long: help.Long(`
			Create a namespace inside an existing index.

			Provide the index name and namespace name. Optionally supply a metadata schema to control which metadata fields are indexed.
		`),
		Example: help.Examples(`
			# create a namespace in an index
			pc index namespace create --index-name "my-index" --name "tenant-a"

			# create a namespace with metadata schema and JSON output
			pc index namespace create --index-name "my-index" --name "tenant-b" --schema "category:keyword" --json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(cmd.Context())

			if strings.TrimSpace(options.indexName) == "" {
				msg.FailMsg("Failed to create namespace: --index-name is required")
				exit.ErrorMsg("Failed to create namespace: --index-name is required")
			}

			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
			if err != nil {
				msg.FailMsg("Failed to create namespace: %s\n", err)
				exit.Error(err, "Failed to create namespace")
			}

			err = runCreateNamespaceCmd(cmd.Context(), ic, options)
			if err != nil {
				msg.FailMsg("Failed to create namespace: %s\n", err)
				exit.Error(err, "Failed to create namespace")
			}
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

func runCreateNamespaceCmd(ctx context.Context, ic NamespaceService, options createNamespaceCmdOptions) error {
	if strings.TrimSpace(options.name) == "" {
		return pcio.Errorf("--name is required")
	}

	req := &pinecone.CreateNamespaceParams{
		Name:   options.name,
		Schema: sdk.BuildMetadataSchema(options.metadataSchema),
	}
	ns, err := ic.CreateNamespace(ctx, req)
	if err != nil {
		return err
	}

	if options.json {
		json := text.IndentJSON(ns)
		pcio.Println(json)
	} else {
		msg.SuccessMsg("Namespace %s created successfully.", options.name)
		presenters.PrintDescribeNamespaceTable(ns)
	}

	return nil
}
