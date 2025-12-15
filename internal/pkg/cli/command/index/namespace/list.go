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

type listNamespaceCmdOptions struct {
	indexName       string
	paginationToken string
	limit           uint32
	prefix          string
	json            bool
}

func NewListNamespaceCmd() *cobra.Command {
	options := listNamespaceCmdOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List namespaces from an index",
		Long:    help.Long(``),
		Example: help.Examples(``),
		Run: func(cmd *cobra.Command, args []string) {
			runListNamespaceCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to list namespaces from")
	cmd.Flags().StringVarP(&options.paginationToken, "pagination-token", "p", "", "pagination token to continue a previous listing operation")
	cmd.Flags().Uint32VarP(&options.limit, "limit", "l", 0, "maximum number of namespaces to list")
	cmd.Flags().StringVar(&options.prefix, "prefix", "", "prefix to filter namespaces by")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runListNamespaceCmd(ctx context.Context, options listNamespaceCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
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
	var prefix *string
	if options.prefix != "" {
		prefix = &options.prefix
	}

	resp, err := ic.ListNamespaces(ctx, &pinecone.ListNamespacesParams{
		PaginationToken: paginationToken,
		Limit:           limit,
		Prefix:          prefix,
	})
	if err != nil {
		msg.FailMsg("Failed to list namespaces: %s", err)
		exit.Error(err, "Failed to list namespaces")
	}

	if options.json {
		json := text.IndentJSON(resp)
		pcio.Println(json)
	} else {
		printTable(resp)
	}
}

func printTable(resp *pinecone.ListNamespacesResponse) {
	writer := presenters.NewTabWriter()
	if resp == nil {
		presenters.PrintEmptyState(writer, "namespaces")
		return
	}

	// Response info
	pcio.Fprintf(writer, "Total Count: %d\n", resp.TotalCount)
	pgToken := "<none>"
	if resp.Pagination != nil {
		pgToken = resp.Pagination.Next
	}
	pcio.Fprintf(writer, "Next Pagination Token: %s\n", pgToken)
	pcio.Fprintf(writer, "\n")

	// Namespaces table
	columns := []string{"NAME", "RECORD COUNT", "INDEXED FIELDS", "SCHEMA"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, ns := range resp.Namespaces {
		schema := "<none>"
		if ns.Schema != nil {
			schema = text.InlineJSON(ns.Schema)
		}
		indexedFields := "<none>"
		if ns.IndexedFields != nil {
			indexedFields = text.InlineJSON(ns.IndexedFields)
		}
		pcio.Fprintf(writer, "%s\t%d\t%s\t%s\n", ns.Name, ns.RecordCount, indexedFields, schema)
	}
	writer.Flush()
}
