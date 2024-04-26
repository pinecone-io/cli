package target

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

var targetHelpTemplate string = `Many API calls take place in the context of a specific project. 
When using the CLI interactively (i.e. via the device authorization flow) you
should use this command to set the current project context for the CLI.

If you're not sure what values to pass to this command, you can discover available 
projects and organizations by running %s.

For automation use cases relying on API-Keys for authentication, there's no need
to specify a project context as the API-Key is already associated with a specific
project in the backend.
`
var targetHelp = fmt.Sprintf(targetHelpTemplate, style.Code("pinecone project list"))

type TargetOptions struct {
	Org     string
	Project string
}

func NewTargetCmd() *cobra.Command {
	options := TargetOptions{}

	cmd := &cobra.Command{
		Use:   "target <command>",
		Short: "Set context for the CLI",
		Long:  targetHelp,
		Run: func(cmd *cobra.Command, args []string) {
			state.TargetOrgName.Set(options.Org)
			state.TargetProjectName.Set(options.Project)

			fmt.Println("✅ Target context updated")
			fmt.Println()
			presenters.PrintTargetContext(state.GetTargetContext())
		},
	}

	// Required options
	cmd.Flags().StringVarP(&options.Org, "org", "o", "", "Organization name")
	cmd.MarkFlagRequired("org")
	cmd.Flags().StringVarP(&options.Project, "project", "p", "", "Project name")
	cmd.MarkFlagRequired("project")

	return cmd
}
