package auth

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type listLocalKeysCmdOptions struct {
	reveal bool
	json   bool
}

var (
	listHelp = help.Long(`
		List the project API keys that the CLI is currently managing in local state.

		The CLI stores only one API key per project, as needed. API key values are
		obscured in the output by default.

		See: https://docs.pinecone.io/reference/cli/authentication
	`)
)

func NewListLocalKeysCmd() *cobra.Command {
	options := listLocalKeysCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the project API keys that the CLI is currently managing in local state",
		Long:  listHelp,
		Example: help.Examples(`
			pc auth local-keys list --reveal
		`),
		Run: func(cmd *cobra.Command, args []string) {
			managedKeys := secrets.GetManagedProjectKeys()
			if options.json {
				maskedMap := maskForJSON(managedKeys, options.reveal)
				json := text.IndentJSON(maskedMap)
				pcio.Println(json)
			} else {
				printTable(managedKeys, options.reveal)
			}
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the API key values in the output")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func printTable(managedKeys map[string]secrets.ManagedKey, reveal bool) {
	writer := presenters.NewTabWriter()

	columns := []string{"PROJECT ID", "API KEY NAME", "API KEY ID", "API KEY VALUE", "ORIGIN", "ORGANIZATION ID"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for projectId, managedKey := range managedKeys {
		keyValue := managedKey.Value
		if !reveal {
			keyValue = presenters.MaskHeadTail(keyValue, 4, 4)
		}
		values := []string{projectId, managedKey.Name, managedKey.Id, keyValue, string(managedKey.Origin), managedKey.OrganizationId}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}

	writer.Flush()
}

func maskForJSON(src map[string]secrets.ManagedKey, reveal bool) map[string]secrets.ManagedKey {
	out := make(map[string]secrets.ManagedKey)
	for projectId, managedKey := range src {
		if !reveal {
			managedKey.Value = presenters.MaskHeadTail(managedKey.Value, 4, 4)
		}
		out[projectId] = managedKey
	}
	return out
}
