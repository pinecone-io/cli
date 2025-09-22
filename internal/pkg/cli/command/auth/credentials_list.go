package auth

import (
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type ListCredentialsCmdOptions struct {
	reveal bool
	json   bool
}

func NewListCredentialsCmd() *cobra.Command {
	options := ListCredentialsCmdOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List the project credentials the CLI is currently managing",
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			managedKeys := secrets.GetManagedProjectKeys()
			if options.json {
				json := text.IndentJSON(managedKeys)
				pcio.Println(json)
			} else {
				printTable(managedKeys, options.reveal)
			}
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the credential key values in the output")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func printTable(managedKeys map[string]secrets.ManagedKey, reveal bool) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"PROJECT ID", "API KEY NAME", "API KEY VALUE", "ORIGIN", "ORGANIZATION ID"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for projectId, managedKey := range managedKeys {
		keyValue := managedKey.Value
		if !reveal {
			keyValue = presenters.MaskHeadTail(keyValue, 4, 4)
		}
		values := []string{projectId, managedKey.Name, keyValue, string(managedKey.Origin), managedKey.OrganizationId}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}

	writer.Flush()
}
