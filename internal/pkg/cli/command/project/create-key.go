package project

import (
	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type CreateApiKeyOptions struct {
	json   bool
	reveal bool
	name   string
}

func NewCreateApiKeyCmd() *cobra.Command {
	options := CreateApiKeyOptions{}

	cmd := &cobra.Command{
		Use:     "create-key",
		Short:   "create an API key in a project",
		GroupID: help.GROUP_PROJECTS_API_KEYS.ID,
		Example: help.Examples([]string{
			"pinecone target -o \"my-org\" -p \"my-project\"",
			"pinecone create-key -n \"my-key\" --reveal",
		}),
		Run: func(cmd *cobra.Command, args []string) {
			projId, err := getTargetProjectId()
			if err != nil {
				msg.FailMsg("No target project set. Use %s to set the target project.", style.Code("pinecone target -o <org> -p <project>"))
				exit.ErrorMsg("No project context set")
			}

			if options.name == "" {
				msg.FailMsg("Name of the key is required")
				exit.ErrorMsg("Name of the key is required")
			}

			existingKeys, err := dashboard.GetApiKeysById(projId)
			if err != nil {
				msg.FailMsg("Failed to list keys: %s", err)
				exit.Error(err)
			}
			for _, key := range existingKeys.Keys {
				if key.UserLabel == options.name {
					msg.FailMsg("Key with name %s already exists", style.Emphasis(options.name))
					msg.HintMsg("See existing keys with %s", style.Code("pinecone project list-keys"))
					exit.ErrorMsg(pcio.Sprintf("Key with name %s already exists", style.Emphasis(options.name)))
				}
			}

			keysResponse, err := dashboard.CreateApiKey(projId, options.name)
			if err != nil {
				msg.FailMsg("Failed to create key: %s", err)
				exit.Error(err)
			}

			var keysToShow []dashboard.Key = []dashboard.Key{}
			if !options.reveal {
				show := dashboard.Key{
					Id:        keysResponse.Key.Id,
					UserLabel: keysResponse.Key.UserLabel,
					Value:     "REDACTED",
				}
				keysToShow = append(keysToShow, show)
			} else {
				keysToShow = append(keysToShow, keysResponse.Key)
			}

			if options.json {
				presentedKey := presentKey(keysToShow[0])
				text.PrettyPrintJSON(presentedKey)
			} else {
				msg.SuccessMsg("Key %s created\n", keysResponse.Key.UserLabel)
				if !options.reveal {
					msg.HintMsg("Run %s to see the key value\n", style.Code("pinecone project list-keys --reveal"))
				}

				printKeysTable(keysToShow)
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "reveal secret key values")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the key to create")
	return cmd
}
