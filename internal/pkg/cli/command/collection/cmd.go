package collection

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

var collectionHelpText = text.WordWrap(`
A collection is a static copy of a pod-based index. It is a non-queryable 
representation of a set of vectors and metadata. You can create a 
collection from a pod-based index, and you can create a new pod-based 
index from a collection. This new index can differ from the original source 
index: the new index can have a different number of pods, a different pod 
type, or a different similarity metric.
`, 80)

func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection <command>",
		Short: "Work with collections",
		Long:  collectionHelpText,
		Example: heredoc.Doc(`
			$ pc collection list
			$ pc collection create --name my-collection --source my-pod-index
			$ pc collection describe my-collection
			$ pc collection delete my-collection
		`),
		GroupID: help.GROUP_VECTORDB.ID,
	}

	cmd.AddCommand(NewCreateCollectionCmd())
	cmd.AddCommand(NewListCollectionsCmd())
	cmd.AddCommand(NewDescribeCollectionCmd())
	cmd.AddCommand(NewDeleteCollectionCmd())

	return cmd
}
