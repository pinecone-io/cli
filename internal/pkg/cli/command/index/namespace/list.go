package namespace

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
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
		Use:   "list",
		Short: "List namespaces from an index",
		Long: help.Long(`
			List namespaces within an index.

			Provide the index name and optionally filter by prefix or paginate with a token and limit. 
			
			Use --json to see the full response including pagination details.
		`),
		Example: help.Examples(`
			# list namespaces for an index
			pc index namespace list --index-name "my-index"

			# list namespaces with a prefix filter and limit
			pc index namespace list --index-name "my-index" --prefix "tenant-" --limit 10

			# continue listing with a pagination token and output JSON
			pc index namespace list --index-name "my-index" --pagination-token "token" --json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			if strings.TrimSpace(options.indexName) == "" {
				msg.FailJSON(options.json, "Failed to list namespaces: --index-name is required")
				exit.ErrorMsg("Failed to list namespaces: --index-name is required")
			}

			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
			if err != nil {
				msg.FailJSON(options.json, "Failed to list namespaces: %s\n", err)
				exit.Error(err, "Failed to list namespaces")
			}

			err = runListNamespaceCmd(ctx, ic, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to list namespaces: %s", err)
				exit.Error(err, "Failed to list namespaces")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to list namespaces from")
	cmd.Flags().StringVarP(&options.paginationToken, "pagination-token", "p", "", "pagination token to continue a previous listing operation")
	cmd.Flags().Uint32VarP(&options.limit, "limit", "l", 0, "maximum number of namespaces to list")
	cmd.Flags().StringVar(&options.prefix, "prefix", "", "prefix to filter namespaces by")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runListNamespaceCmd(ctx context.Context, ic NamespaceService, options listNamespaceCmdOptions) error {
	if strings.TrimSpace(options.indexName) == "" {
		return fmt.Errorf("--index-name is required")
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
		return err
	}

	if options.json {
		json := text.IndentJSON(resp)
		fmt.Fprintln(os.Stdout, json)
	} else {
		printTable(resp)
	}

	return nil
}

func printTable(resp *pinecone.ListNamespacesResponse) {
	writer := presenters.NewTabWriter()
	if resp == nil {
		presenters.PrintEmptyState(writer, "namespaces")
		return
	}

	// Response info
	fmt.Fprintf(writer, "Total Count: %d\n", resp.TotalCount)
	pgToken := "<none>"
	if resp.Pagination != nil && resp.Pagination.Next != "" {
		pgToken = resp.Pagination.Next
	}
	fmt.Fprintf(writer, "Next Pagination Token: %s\n", pgToken)
	fmt.Fprintf(writer, "\n")

	// Namespaces table
	columns := []string{"NAME", "RECORD COUNT", "INDEXED FIELDS", "SCHEMA"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	for _, ns := range resp.Namespaces {
		schema := "<none>"
		if ns.Schema != nil {
			schema = text.InlineJSON(ns.Schema)
		}
		indexedFields := "<none>"
		if ns.IndexedFields != nil {
			indexedFields = text.InlineJSON(ns.IndexedFields)
		}
		fmt.Fprintf(writer, "%s\t%d\t%s\t%s\n", ns.Name, ns.RecordCount, indexedFields, schema)
	}
	writer.Flush()
}
