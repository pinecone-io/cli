package apiKey

import (
	"sort"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type ListKeysCmdCmdOptions struct {
	projectId string
	json      bool
}

func NewListKeysCmd() *cobra.Command {
	options := ListKeysCmdCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List API keys for the target project, or a specific project ID",
		Example: help.Examples(`
			# List API keys for the target project
			pc target --org "org-name" --project "project-name"
			pc api-key list

			# List API keys for a specific project
			pc api-key list --id "project-id"
		`),
		GroupID: help.GROUP_API_KEYS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			var err error
			projId := options.projectId
			if projId == "" {
				projId, err = state.GetTargetProjectId()
				if err != nil {
					msg.FailMsg("No target project set, and no project ID provided. Use %s to set the target project. Use %s to create the key in a specific project.", style.Code("pc target -o <org> -p <project>"), style.Code("pc api-key create -i <project-id> -n <name>"))
					exit.ErrorMsg("No project ID provided, and no target project set")
				}
			}

			keysResponse, err := ac.APIKey.List(cmd.Context(), projId)
			if err != nil {
				msg.FailMsg("Failed to list API keys: %s", err)
				exit.Error(err)
			}

			// Sort keys alphabetically by name
			sortedKeys := keysResponse
			sort.Slice(sortedKeys, func(i, j int) bool {
				return sortedKeys[i].Name < sortedKeys[j].Name
			})

			if options.json {
				json := text.IndentJSON(sortedKeys)
				pcio.Println(json)
			} else {
				printTable(sortedKeys)
			}
		},
	}

	cmd.Flags().StringVarP(&options.projectId, "id", "i", "", "ID of the project to list the keys for if not the target project")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")
	return cmd
}

func printTable(keys []*pinecone.APIKey) {
	pcio.Printf("Organization: %s (ID: %s)\n", style.Emphasis(state.TargetOrg.Get().Name), style.Emphasis(state.TargetOrg.Get().Id))
	pcio.Printf("Project: %s (ID: %s)\n", style.Emphasis(state.TargetProj.Get().Name), style.Emphasis(state.TargetProj.Get().Id))
	pcio.Println()
	pcio.Println(style.Heading("API Keys"))
	pcio.Println()

	writer := presenters.NewTabWriter()

	columns := []string{"NAME", "ID", "PROJECT ID", "ROLES"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, key := range keys {
		values := []string{
			key.Name,
			key.Id,
			key.ProjectId,
			strings.Join(key.Roles, ", "),
		}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}

	writer.Flush()
}
