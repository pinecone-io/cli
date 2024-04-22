package index

import (
	"fmt"
	"os"
	"context"

	"github.com/spf13/cobra"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
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
			key := os.Getenv("PINECONE_API_KEY")
			fmt.Println("list called with key:", key)

			ctx := context.Background()

			pc, err := pinecone.NewClient(pinecone.NewClientParams{
				ApiKey: key,
			})
		
			if err != nil {
				exit.Error(err)
			}
		
			idxs, err := pc.ListIndexes(ctx)
			if err != nil {
				exit.Error(err)
			}

			fmt.Println(idxs)

			text.PrettyPrintJSON(idxs)
		},
	}

	return cmd
}