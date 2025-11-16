package project

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type deleteProjectCmdOptions struct {
	projectId        string
	skipConfirmation bool
	json             bool
}

func NewDeleteProjectCmd() *cobra.Command {
	options := deleteProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a project (the target project, or a specific project ID)",
		Example: help.Examples(`
			# Delete the target project
			pc project delete

			# Delete a specific project by ID
			pc project delete --id "project-id"
		`),
		GroupID: help.GROUP_PROJECTS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ac := sdk.NewPineconeAdminClient(ctx)

			projId := options.projectId
			var err error
			if projId == "" {
				projId, err = state.GetTargetProjectId()
				if err != nil {
					msg.FailMsg("No target project set and no project ID provided. Use %s to set the target project. Use %s to delete a specific project.", style.Code("pc target -p <project>"), style.Code("pc project delete -i <project-id>"))
					exit.ErrorMsg("No project ID provided, and no target project set")
				}
			}

			projToDelete, err := ac.Project.Describe(ctx, projId)
			if err != nil {
				msg.FailMsg("Failed to retrieve project information: %s\n", err)
				msg.HintMsg("To see a list of projects in the organization, run %s", style.Code("pc project list"))
				exit.Error(err, "Failed to retrieve project information")
			}

			verifyNoIndexes(ctx, projToDelete.Id, projToDelete.Name)
			verifyNoCollections(ctx, projToDelete.Id, projToDelete.Name)

			if !options.skipConfirmation {
				confirmDelete(projToDelete.Name)
			}

			err = ac.Project.Delete(ctx, projToDelete.Id)
			if err != nil {
				msg.FailMsg("Failed to delete project %s: %s\n", style.Emphasis(projToDelete.Name), err)
				exit.Errorf(err, "Failed to delete project %s", style.Emphasis(projToDelete.Name))
			}

			// Clear target project if the deleted project is the target project
			if state.TargetProj.Get().Name == projToDelete.Name {
				state.TargetProj.Set(state.TargetProject{
					Id:   "",
					Name: "",
				})
			}
			msg.SuccessMsg("Project %s deleted.\n", style.Emphasis(projToDelete.Name))
		},
	}

	// optional flags
	cmd.Flags().StringVarP(&options.projectId, "id", "i", "", "ID of the project to delete")
	cmd.Flags().BoolVar(&options.skipConfirmation, "skip-confirmation", false, "Skip the deletion confirmation prompt")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func confirmDelete(projectName string) {
	msg.WarnMsg("This will delete the project %s in organization %s.", style.Emphasis(projectName), style.Emphasis(state.TargetOrg.Get().Name))
	msg.WarnMsg("This action cannot be undone.")

	// Prompt the user
	fmt.Print("Do you want to continue? (y/N): ")

	// Read the user's input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		pcio.Println(fmt.Errorf("Error reading input: %w", err))
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

func verifyNoIndexes(ctx context.Context, projectId string, projectName string) {
	// Check if project contains indexes
	pc := sdk.NewPineconeClientForProjectById(ctx, projectId)

	idxs, err := pc.ListIndexes(ctx)
	if err != nil {
		msg.FailMsg("Failed to list indexes: %s\n", err)
		exit.Error(err, "Failed to list indexes")
	}
	if len(idxs) > 0 {
		msg.FailMsg("Project %s contains indexes. Delete the indexes before deleting the project.", style.Emphasis(projectName))
		msg.HintMsg("To see indexes in this project, run %s", pcio.Sprintf(style.Code("pc target -p \"%s\" && pc index list"), projectName))
		exit.ErrorMsgf("Project %s contains indexes. Delete the indexes before deleting the project.", style.Emphasis(projectName))
	}
}

func verifyNoCollections(ctx context.Context, projectId string, projectName string) {
	// Check if project contains collections
	pc := sdk.NewPineconeClientForProjectById(ctx, projectId)

	collections, err := pc.ListCollections(ctx)
	if err != nil {
		msg.FailMsg("Failed to list collections: %s\n", err)
		exit.Errorf(err, "Failed to list collections")
	}
	if len(collections) > 0 {
		msg.FailMsg("Project %s contains collections. Delete the collections before deleting the project.", style.Emphasis(projectName))
		msg.HintMsg("To see collections in this project, run %s", pcio.Sprintf(style.Code("pc target -p \"%s\" && pc collection list"), projectName))
		exit.ErrorMsgf("Project %s contains collections. Delete the collections before deleting the project.", style.Emphasis(projectName))
	}
}
