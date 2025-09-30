package index

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

var helpText = text.WordWrap(`An index is the highest-level organizational unit of 
vector data in Pinecone. It accepts and stores vectors, serves queries 
over the vectors it contains, and does other vector operations over 
its contents.`, 80)

func NewIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Work with indexes",
		Long:  helpText,
		Example: help.Examples(`
			$ pc index list
			$ pc index create --name my-index --dimension 1536 --metric cosine --cloud aws --region us-east-1
			$ pc index describe --name my-index
			$ pc index delete --name my-index
		`),
		GroupID: help.GROUP_VECTORDB.ID,
	}

	cmd.AddCommand(NewDescribeCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateIndexCmd())
	cmd.AddCommand(NewCreateServerlessCmd())
	cmd.AddCommand(NewCreatePodCmd())
	cmd.AddCommand(NewConfigureIndexCmd())
	cmd.AddCommand(NewDeleteCmd())

	return cmd
}
