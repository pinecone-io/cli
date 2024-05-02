package project

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteProjectCmdOptions struct {
	name string
	json bool
	yes  bool
}

func NewDeleteProjectCmd() *cobra.Command {
	options := DeleteProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a project in the target org",
		Example: help.Examples([]string{
			"pinecone target -o \"my-org\"",
			"pinecone project delete --name=\"demo\"",
			"pinecone project delete --name=\"demo\" --yes",
		}),
		GroupID: help.GROUP_PROJECTS_CRUD.ID,
		Run: func(cmd *cobra.Command, args []string) {
			orgId, err := getTargetOrgId()
			orgName := state.TargetOrg.Get().Name
			if err != nil {
				msg.FailMsg("No target organization set. Use %s to set the organization context.", style.Code("pinecone target -o <org>"))
				cmd.Help()
				exit.ErrorMsg("No organization context set")
			}

			projToDelete, err := dashboard.GetProjectByName(orgName, options.name)
			if err != nil {
				msg.FailMsg("Failed to retrieve project information: %s\n", err)
				msg.HintMsg("To see a list of projects in the organization, run %s", style.Code("pinecone project list"))
				exit.Error(err)
			}

			verifyNoIndexes(orgName, projToDelete.Name)
			verifyNoCollections(orgName, projToDelete.Name)

			if !options.yes {
				confirmDelete(options.name, orgId)
			}

			resp, err := dashboard.DeleteProject(orgId, projToDelete.Id)
			if err != nil {
				msg.FailMsg("Failed to delete project %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}
			if !resp.Success {
				msg.FailMsg("Failed to delete project %s\n", style.Emphasis(options.name))
				exit.Error(pcio.Errorf("Delete project %s call returned 200 but with success=false in the body", options.name))
			}

			// Clear target project if the deleted project is the target project
			if state.TargetProj.Get().Name == options.name {
				state.TargetProj.Set(&state.TargetProject{
					Id:   "",
					Name: "",
				})
			}
			msg.SuccessMsg("Project %s deleted.\n", style.Emphasis(options.name))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().BoolVar(&options.yes, "yes", false, "skip confirmation prompt")
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

func getTargetProjectId() (string, error) {
	projId := state.TargetProj.Get().Id
	if projId == "" {
		return "", pcio.Errorf("no target project set")
	}
	return projId, nil
}

func confirmDelete(projectName, orgId string) {
	msg.WarnMsg("This will delete the project %s in organization %s.", style.Emphasis(projectName), style.Emphasis(state.TargetOrg.Get().Name))
	msg.WarnMsg("This action cannot be undone.")

	// Prompt the user
	fmt.Print("Do you want to continue? (y/N): ")

	// Read the user's input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Trim any whitespace from the input and convert to lowercase
	input = strings.TrimSpace(strings.ToLower(input))

	// Check if the user entered "y" or "yes"
	if input == "y" || input == "yes" {
		msg.InfoMsg("You chose to continue delete.")
	} else {
		msg.InfoMsg("Operation canceled.")
		exit.Success()
	}
}

func verifyNoIndexes(orgName string, projectName string) {
	// Check if project contains indexes
	pc := sdk.NewPineconeClientForUserProjectByName(orgName, projectName)
	ctx := context.Background()

	idxs, err := pc.ListIndexes(ctx)
	if err != nil {
		msg.FailMsg("Failed to list indexes: %s\n", err)
		exit.Error(err)
	}
	if len(idxs) > 0 {
		msg.FailMsg("Project %s contains indexes. Delete the indexes before deleting the project.", style.Emphasis(projectName))
		msg.HintMsg("To see indexes in this project, run %s", pcio.Sprintf(style.Code("pinecone target -p \"%s\" && pinecone index list"), projectName))
		exit.Error(pcio.Errorf("project contains indexes"))
	}
}

func verifyNoCollections(orgName string, projectName string) {
	// Check if project contains collections
	pc := sdk.NewPineconeClientForUserProjectByName(orgName, projectName)
	ctx := context.Background()

	collections, err := pc.ListCollections(ctx)
	if err != nil {
		msg.FailMsg("Failed to list collections: %s\n", err)
		exit.Error(err)
	}
	if len(collections) > 0 {
		msg.FailMsg("Project %s contains collections. Delete the collections before deleting the project.", style.Emphasis(projectName))
		msg.HintMsg("To see collections in this project, run %s", pcio.Sprintf(style.Code("pinecone target -p \"%s\" && pinecone collection list"), projectName))
		exit.Error(pcio.Errorf("project contains collections"))
	}
}
