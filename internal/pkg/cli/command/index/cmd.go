package index

import (
	"github.com/spf13/cobra"
)

var helpText = `An index is the highest-level organizational unit of 
vector data in Pinecone. It accepts and stores vectors, serves queries 
over the vectors it contains, and does other vector operations over 
its contents.`

func NewIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index <command>",
		Short: "Work with indexes",
		Long:  helpText,
	}

	cmd.AddCommand(NewDescribeCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateServerlessCmd())

	return cmd
}
