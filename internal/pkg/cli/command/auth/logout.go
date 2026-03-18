package auth

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type logoutCmdOptions struct {
	json bool
}

func NewLogoutCmd() *cobra.Command {
	options := logoutCmdOptions{}

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Clears locally stored API keys (managed and default), service account details (client ID and secret), and user login token. Also clears organization and project target context.",
		Example: help.Examples(`
			pc auth logout
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			secrets.ConfigFile.Clear()
			state.ConfigFile.Clear()

			if options.json {
				fmt.Println(text.IndentJSON(struct {
					Status string `json:"status"`
				}{Status: "logged_out"}))
				return
			}

			msg.SuccessMsg("API keys and user access tokens cleared.")
			msg.SuccessMsg("State cleared.")
		},
	}

	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output result as JSON")

	return cmd
}
