package collection

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection <command>",
		Short: "Work with collections",
		Long: help.Long(`
			A collection is a static copy of an index. It is a non-queryable 
			representation of a set of vectors and metadata. You can create a 
			collection from an index, and you can create a new index from a 
			collection. 
			
			This new index can differ from the original source index: the 
			new index can have a different number of pods, a different pod type, 
			or a different similarity metric.
		`),
		GroupID: help.GROUP_VECTORDB.ID,
	}

	cmd.AddCommand(NewCreateCollectionCmd())
	cmd.AddCommand(NewListCollectionsCmd())
	cmd.AddCommand(NewDescribeCollectionCmd())
	cmd.AddCommand(NewDeleteCollectionCmd())

	return cmd
}
