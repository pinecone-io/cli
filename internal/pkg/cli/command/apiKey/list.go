package apiKey

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
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
		Short: "List the API keys in a specific project by ID or the target project",
		Example: heredoc.Doc(`
		$ pc target -o "my-org" -p "my-project"
		$ pc api-key list
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
				fmt.Println(json)
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
	fmt.Printf("Organization: %s (ID: %s)\n", style.Emphasis(state.TargetOrg.Get().Name), style.Emphasis(state.TargetOrg.Get().Id))
	fmt.Printf("Project: %s (ID: %s)\n", style.Emphasis(state.TargetProj.Get().Name), style.Emphasis(state.TargetProj.Get().Id))
	fmt.Println()
	fmt.Println(style.Heading("API Keys"))
	fmt.Println()

	writer := presenters.NewTabWriter()

	columns := []string{"NAME", "ID", "PROJECT ID", "ROLES"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	for _, key := range keys {
		values := []string{
			key.Name,
			key.Id,
			key.ProjectId,
			strings.Join(key.Roles, ", "),
		}
		fmt.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}

	writer.Flush()
}
