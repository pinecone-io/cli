package collection

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/docslinks"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var collectionHelpText = fmt.Sprintf(`To learn more about collections, please see %s`, docslinks.UnderstandingCollectionsGuide)

func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection <command>",
		Short:   "Work with collections",
		Long:    collectionHelpText,
		GroupID: help.GROUP_VECTORDB.ID,
	}

	cmd.AddCommand(NewCreateCollectionCmd())
	cmd.AddCommand(NewListCollectionsCmd())
	cmd.AddCommand(NewDescribeCollectionCmd())
	cmd.AddCommand(NewDeleteCollectionCmd())

	return cmd
}
