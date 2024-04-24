package collection

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/docslinks"
	"github.com/spf13/cobra"
)

var collectionHelpText = fmt.Sprintf(`To learn more about collections, please see %s`, docslinks.UnderstandingCollectionsGuide)

func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection <command>",
		Short: "Work with collections",
		Long:  collectionHelpText,
	}

	cmd.AddCommand(NewListCollectionsCmd())
	cmd.AddCommand(NewCreateCollectionCmd())

	return cmd
}
