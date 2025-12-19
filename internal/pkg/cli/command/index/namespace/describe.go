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
		Use:   "describe",
		Short: "Describe a namespace from an index by name",
		Long: help.Long(`
			Describe a namespace by name, including record counts and schema configuration.
		`),
		Example: help.Examples(`
			# describe a namespace
			pc index namespace describe --index-name "my-index" --name "tenant-a"

			# describe a namespace and return JSON
			pc index namespace describe --index-name "my-index" --name "tenant-a" --json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			if strings.TrimSpace(options.indexName) == "" {
				msg.FailMsg("Failed to describe namespace: --index-name is required")
				exit.ErrorMsg("Failed to describe namespace: --index-name is required")
			}

			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
			if err != nil {
				msg.FailMsg("Failed to describe namespace: %s\n", err)
				exit.Error(err, "Failed to describe namespace")
			}

			err = runDescribeNamespaceCmd(ctx, ic, options)
			if err != nil {
				msg.FailMsg("Failed to describe namespace: %s", err)
				exit.Error(err, "Failed to describe namespace")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to describe the namespace from")
	cmd.Flags().StringVar(&options.name, "name", "", "name of the namespace to describe")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runDescribeNamespaceCmd(ctx context.Context, ic NamespaceService, options describeNamespaceCmdOptions) error {
	if strings.TrimSpace(options.name) == "" {
		return pcio.Errorf("--name is required")
	}

	ns, err := ic.DescribeNamespace(ctx, options.name)
	if err != nil {
		return err
	}

	if options.json {
		json := text.IndentJSON(ns)
		pcio.Println(json)
	} else {
		presenters.PrintDescribeNamespaceTable(ns)
	}

	return nil
}
