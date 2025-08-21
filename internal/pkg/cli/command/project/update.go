package project

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type UpdateProjectCmdOptions struct {
	projectId               string
	name                    string
	forceEncryptionWithCMEK bool
	maxPods                 int

	json bool
}

func NewUpdateProjectCmd() *cobra.Command {
	options := UpdateProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing project with the specified configuration",
		Example: heredoc.Doc(`
		$ pc project update --id <project-id> --name <new-name> --max-pods <new-max-pods>
		`),
		GroupID: help.GROUP_PROJECTS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()
			if len(args) == 0 {
				msg.FailMsg("No project ID provided, please provide a project ID")
				exit.ErrorMsg("No project ID provided")
			}

			projectId := args[0]

			project, err := ac.Project.Update(context.Background(), projectId, &pinecone.UpdateProjectParams{
				Name:                    &options.name,
				ForceEncryptionWithCmek: &options.forceEncryptionWithCMEK,
				MaxPods:                 &options.maxPods,
			})
			if err != nil {
				msg.FailMsg("Failed to update project %s: %s\n", projectId, err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(project)
				pcio.Println(json)
				return
			}

			msg.SuccessMsg("Project %s updated successfully.", project.Id)
			presenters.PrintDescribeProjectTable(project)
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.projectId, "id", "i", "", "ID of the project to update")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of project to update")
	cmd.Flags().BoolVarP(&options.forceEncryptionWithCMEK, "force-encryption-with-cmek", "f", false, "force encryption with CMEK")
	cmd.Flags().IntVarP(&options.maxPods, "max-pods", "p", 0, "max pods for the project")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
