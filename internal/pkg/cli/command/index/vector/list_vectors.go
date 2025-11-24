package vector

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

type listVectorsCmdOptions struct {
	indexName       string
	namespace       string
	limit           uint32
	paginationToken string
	json            bool
}

func NewListVectorsCmd() *cobra.Command {
	options := listVectorsCmdOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List vectors in an index",
		Example: help.Examples(`
			pc index vector list --index-name my-index --namespace my-namespace
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runListVectorsCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to list vectors from")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to list vectors from")
	cmd.Flags().Uint32VarP(&options.limit, "limit", "l", 0, "maximum number of vectors to list")
	cmd.Flags().StringVarP(&options.paginationToken, "pagination-token", "p", "", "pagination token to continue a previous listing operation")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runListVectorsCmd(ctx context.Context, options listVectorsCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	// Get IndexConnection
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	var limit *uint32
	if options.limit > 0 {
		limit = &options.limit
	}
	var paginationToken *string
	if options.paginationToken != "" {
		paginationToken = &options.paginationToken
	}

	resp, err := ic.ListVectors(ctx, &pinecone.ListVectorsRequest{
		Limit:           limit,
		PaginationToken: paginationToken,
	})
	if err != nil {
		msg.FailMsg("Failed to list vectors: %s", err)
		exit.Error(err, "Failed to list vectors")
	}

	if options.json {
		json := text.IndentJSON(resp)
		pcio.Println(json)
	} else {
		presenters.PrintListVectorsTable(resp)
	}
}
