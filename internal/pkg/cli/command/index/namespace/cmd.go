package namespace

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	namespaceHelp = help.Long(`
		Work with namespaces in a Pinecone index.

		Use these commands to create, list, describe, and delete namespaces within an index.

		See: https://docs.pinecone.io/guides/manage-data/manage-namespaces
	`)
)

func NewNamespaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "Work with namespaces in an index",
		Long:  namespaceHelp,
		Example: help.Examples(`
			# list namespaces for an index
			pc index namespace list --index-name "my-index"

            # create a new namespace in an index
			pc index namespace create --index-name "my-index" --name "tenant-c"

			# describe a specific namespace
			pc index namespace describe --index-name "my-index" --name "tenant-a"
		`),
		GroupID: help.GROUP_INDEX_NAMESPACE.ID,
	}

	cmd.AddCommand(NewCreateNamespaceCmd())
	cmd.AddCommand(NewListNamespaceCmd())
	cmd.AddCommand(NewDescribeNamespaceCmd())
	cmd.AddCommand(NewDeleteNamespaceCmd())

	return cmd
}
