package project

import (
	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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
		Use:     "create",
		Short:   "create a project in the target org",
		GroupID: help.GROUP_PROJECTS_CRUD.ID,
		Example: help.Examples([]string{
			"pc target -o \"my-org\"",
			"pc project create --name=\"demo\"",
		}),
		Run: func(cmd *cobra.Command, args []string) {
			orgId, err := getTargetOrgId()
			if err != nil {
				msg.FailMsg("No target organization set. Use %s to set the organization context.", style.Code("pc target -o <org>"))
				cmd.Help()
				exit.ErrorMsg("No organization context set")
			}

			proj, err := dashboard.CreateProject(orgId, options.name, options.pod_quota)
			if err != nil {
				msg.FailMsg("Failed to create project %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}
			if !proj.Success {
				msg.FailMsg("Failed to create project %s\n", style.Emphasis(options.name))
				exit.Error(pcio.Errorf("Create project call returned 200 but with success=false in the body%s", options.name))
			}
			msg.SuccessMsg("Project %s created successfully.\n", style.Emphasis(proj.Project.Name))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the project")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int32VarP(&options.pod_quota, "pod_quota", "p", 5, "maximum number of pods allowed in the project across all indexes")
	return cmd
}
