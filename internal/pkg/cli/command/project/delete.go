package project

import (
	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteProjectCmdOptions struct {
	name string
	json bool
}

func NewDeleteProjectCmd() *cobra.Command {
	options := DeleteProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a project in the target org",
		Example: help.Examples([]string{
			"pinecone target -o \"my-org\"",
			"pinecone project delete --name=\"demo\"",
		}),
		Run: func(cmd *cobra.Command, args []string) {
			orgId, err := getTargetOrgId()
			if err != nil {
				msg.FailMsg("No target organization set. Use %s to set the organization context.", style.Code("pinecone target -o <org>"))
				cmd.Help()
				exit.ErrorMsg("No organization context set")
			}

			orgs, err := dashboard.ListOrganizations()
			if err != nil {
				msg.FailMsg("Failed to retrieve org information: %s\n", err)
				exit.Error(err)
			}

			var projectId string
			for _, org := range orgs.Organizations {
				if org.Id == orgId {
					for _, proj := range org.Projects {
						if proj.GlobalProject.Name == options.name {
							projectId = proj.GlobalProject.Id
						}
					}
				}
			}
			if projectId == "" {
				msg.FailMsg("Project %s not found in organization %s. Did you already delete it?\n", style.Emphasis(options.name), style.Emphasis(state.TargetOrg.Get().Name))
				msg.HintMsg("To see a list of projects in the organization, run %s", style.Code("pinecone project list"))
				exit.Error(pcio.Errorf("project not found"))
			}

			resp, err := dashboard.DeleteProject(orgId, projectId)
			if err != nil {
				msg.FailMsg("Failed to delete project %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}
			if !resp.Success {
				msg.FailMsg("Failed to delete project %s: %s\n", style.Emphasis(options.name))
			}
			msg.SuccessMsg("Project %s deleted.\n", style.Emphasis(options.name))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the project")
	cmd.MarkFlagRequired("name")
	return cmd
}

func getTargetOrgId() (string, error) {
	orgId := state.TargetOrg.Get().Id
	if orgId == "" {
		return "", pcio.Errorf("no target organization set")
	}
	return orgId, nil
}
