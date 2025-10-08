package collection

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	collectionHelp = help.Long(`
		Create, describe, list, and delete collections for a pod-based index.

		Collections are static snapshots of pod-based indexes. They preserve vector
		data and metadata but cannot be queried directly. Collections are useful for:

		- Protecting an index from manual or system failures.
		- Temporarily shutting down an index.
		- Copying the data from one index into a different index.
		- Making a backup of your index.
		- Experimenting with different index configurations.

		Collections only work with pod-based indexes (not serverless)

		See: https://docs.pinecone.io/guides/indexes/understanding-collections
	`)
)

func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection <command>",
		Short:   "Work with collections (pod-based indexes only)",
		Long:    collectionHelp,
		GroupID: help.GROUP_VECTORDB.ID,
	}

	cmd.AddCommand(NewCreateCollectionCmd())
	cmd.AddCommand(NewListCollectionsCmd())
	cmd.AddCommand(NewDescribeCollectionCmd())
	cmd.AddCommand(NewDeleteCollectionCmd())

	return cmd
}
