package index

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	indexHelp = help.Long(`
		Work with Pinecone indexes.

		An index is the highest-level organizational unit of vector data in Pinecone. 
		You store vector data in an index, and perform queries and operations on the data stored.
		There are two types of indexes: dense and sparse.
		
		See: https://docs.pinecone.io/guides/index-data/indexing-overview
	`)
)

func NewIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Work with indexes",
		Long:  indexHelp,
		Example: help.Examples(`
			pc index list
			pc index create --name my-index --dimension 1536 --metric cosine --cloud aws --region us-east-1
			pc index describe --name my-index
			pc index delete --name my-index
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
