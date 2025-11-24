package vector

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	vectorHelp = help.Long(`
		Work with vectors (records) in a Pinecone index.

		Use these commands to upsert, fetch, list, update, delete, and query data
		within an index. All commands require --index-name and may optionally target
		a --namespace.

		See: https://docs.pinecone.io/guides/index-data/data-ingestion-overview
	`)
)

func NewVectorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vector",
		Aliases: []string{"vectors", "record", "records"},
		Short:   "Work with data in an index",
		Long:    vectorHelp,
		Example: help.Examples(`
			pc index vector list-vectors --index-name my-index --namespace my-namespace
			pc index vector fetch --index-name my-index --ids '["123","456"]'
			pc index vector upsert --index-name my-index ./vectors.json
			pc index vector update --index-name my-index --id doc-123 --metadata '{"genre":"sci-fi"}'
			pc index vector delete-vectors --index-name my-index --ids doc-123
			pc index vector query --index-name my-index --vector @./vector.json --top-k 10
			pc index vector describe-stats --index-name my-index --filter '{"genre":"sci-fi"}'
		`),
		GroupID: help.GROUP_INDEX_DATA.ID,
	}

	cmd.AddCommand(NewUpsertCmd())
	cmd.AddCommand(NewFetchCmd())
	cmd.AddCommand(NewQueryCmd())
	cmd.AddCommand(NewListVectorsCmd())
	cmd.AddCommand(NewDeleteVectorsCmd())
	cmd.AddCommand(NewUpdateCmd())

	return cmd
}
