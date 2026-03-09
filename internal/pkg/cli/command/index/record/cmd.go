package record

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewRecordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "Work with text records in an integrated index",
		Long: help.Long(`
			Work with text records in an integrated Pinecone index.

			Use these commands to upsert raw text records and run semantic search against
			them. All commands require --index-name and may optionally target a
			--namespace.
		`),
		Example: help.Examples(`
			pc index record upsert --index-name my-index --namespace my-namespace --body ./records.jsonl
			pc index record search --index-name my-index --namespace my-namespace --inputs '{"text":"search query"}'
			pc index record search --index-name my-index --namespace my-namespace --body ./search.json
		`),
		GroupID: help.GROUP_INDEX_DATA.ID,
	}

	cmd.AddCommand(NewUpsertCmd())
	cmd.AddCommand(NewSearchCmd())

	return cmd
}
