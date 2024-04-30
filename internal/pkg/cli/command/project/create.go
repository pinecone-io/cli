package project

import (
	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type CreateProjectCmdOptions struct {
	name      string
	pod_quota int32
	json      bool
}

func NewCreateProjectCmd() *cobra.Command {
	options := CreateProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a project in the target org",
		Example: help.Examples([]string{
			"pinecone target -o \"my-org\"",
			"pinecone project create --name=\"demo\"",
		}),
		Run: func(cmd *cobra.Command, args []string) {
			orgId, err := getTargetOrgId()
			if err != nil {
				msg.FailMsg("No target organization set. Use %s to set the organization context.", style.Code("pinecone target -o <org>"))
				cmd.Help()
				exit.ErrorMsg("No organization context set")
			}

			proj, err := dashboard.CreateProject(orgId, options.name, options.pod_quota)
			if err != nil {
				msg.FailMsg("Failed to create project %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}
			if !proj.Success {
				msg.FailMsg("Failed to create project %s: %s\n", style.Emphasis(options.name))
			}
			msg.SuccessMsg("Project %s created successfully.\n", style.Emphasis(proj.GlobalProject.Name))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the project")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int32VarP(&options.pod_quota, "pod_quota", "p", 5, "maximum number of pods allowed in the project across all indexes")
	return cmd
}
