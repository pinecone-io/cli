package index

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/client"
)

var listHelpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "See the list of indexes in your project",
		Long: listHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			pc := client.NewPineconeClient()
			ctx := context.Background()
		
			idxs, err := pc.ListIndexes(ctx)
			if err != nil {
				exit.Error(err)
			}

			text.PrettyPrintJSON(idxs)
		},
	}

	return cmd
}