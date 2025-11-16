package index

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

type fetchCmdOptions struct {
	name      string
	namespace string
	ids       []string
	json      bool
}

func NewFetchCmd() *cobra.Command {
	options := fetchCmdOptions{}
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch vectors by ID from an index",
		Example: help.Examples(`
			pc index fetch --name my-index --ids 123, 456, 789
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runFetchCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringSliceVarP(&options.ids, "ids", "i", []string{}, "IDs of vectors to fetch")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the index to fetch from")
	cmd.Flags().StringVar(&options.namespace, "namespace", "", "namespace to fetch from")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("ids")

	return cmd
}

func runFetchCmd(ctx context.Context, options fetchCmdOptions) {
	pc := sdk.NewPineconeClient()

	// Default namespace
	ns := options.namespace
	if options.namespace != "" {
		ns = options.namespace
	}
	if ns == "" {
		ns = "__default__"
	}

	ic, err := sdk.NewIndexConnection(ctx, pc, options.name, ns)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	vectors, err := ic.FetchVectors(ctx, options.ids)
	if err != nil {
		exit.Error(err, "Failed to fetch vectors")
	}

	if options.json {
		json := text.IndentJSON(vectors)
		pcio.Println(json)
	} else {
		presenters.PrintFetchVectorsTable(vectors)
	}
}
