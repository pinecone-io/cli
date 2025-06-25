package project

import (
	"sort"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type ListKeysCmdCmdOptions struct {
	json   bool
	reveal bool
}

func NewListKeysCmd() *cobra.Command {
	options := ListKeysCmdCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list-keys",
		Short: "list the API keys in a project",
		Example: help.Examples([]string{
			"pinecone target -o \"my-org\" -p \"my-project\"",
			"pinecone list-keys --reveal",
		}),
		GroupID: help.GROUP_PROJECTS_API_KEYS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			projId, err := getTargetProjectId()
			if err != nil {
				msg.FailMsg("No target project set. Use %s to set the target project.", style.Code("pinecone target -o <org> -p <project>"))
				exit.ErrorMsg("No project context set")
			}

			keysResponse, err := dashboard.GetApiKeysById(projId)
			if err != nil {
				msg.FailMsg("Failed to list keys: %s", err)
				exit.Error(err)
			}

			// Sort keys alphabetically
			sortedKeys := keysResponse.Keys
			sort.Slice(sortedKeys, func(i, j int) bool {
				return sortedKeys[i].UserLabel < sortedKeys[j].UserLabel
			})

			// Unless --reveal, redact secret key values
			var keysToShow []dashboard.Key = []dashboard.Key{}
			if !options.reveal {
				for _, key := range sortedKeys {
					keysToShow = append(keysToShow, dashboard.Key{
						Id:        key.Id,
						UserLabel: key.UserLabel,
						Value:     "REDACTED",
					})
				}
			} else {
				keysToShow = sortedKeys
			}

			// Display output
			if options.json {
				presentedKeys := []PresentedKey{}
				for _, key := range keysToShow {
					presentedKeys = append(presentedKeys, presentKey(key))
				}
				json := text.IndentJSON(presentedKeys)
				pcio.Println(json)
			} else {
				pcio.Printf("org: %s\n", style.Emphasis(state.TargetOrg.Get().Name))
				pcio.Printf("project: %s\n", style.Emphasis(state.TargetProj.Get().Name))
				pcio.Println()
				pcio.Println(style.Heading("API Keys"))
				if !options.reveal {
					msg.HintMsg("To see the key values, add the %s flag", style.Code("--reveal"))
				}

				pcio.Println()
				printKeysTable(keysToShow)

			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "reveal secret key values")
	return cmd
}

func printKeysTable(keys []dashboard.Key) {
	w := presenters.NewTabWriter()

	columns := []string{"NAME", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(w, header)

	for _, key := range keys {
		var values []string
		if key.Value == "REDACTED" {
			values = []string{key.UserLabel, style.StatusRed(key.Value)}
		} else {
			values = []string{key.UserLabel, key.Value}
		}
		pcio.Fprintf(w, strings.Join(values, "\t")+"\n")
	}

	w.Flush()
}

type PresentedKey struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Id    string `json:"id"`
}

func presentKey(key dashboard.Key) PresentedKey {
	return PresentedKey{
		Name:  key.UserLabel,
		Value: key.Value,
		Id:    key.Id,
	}
}
