package project

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeProjectCmdOptions struct {
	projectId string
	json      bool
}

func NewDescribeProjectCmd() *cobra.Command {
	options := DescribeProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Describe a specific project by ID or the target project",
		GroupID: help.GROUP_PROJECTS.ID,
		Example: help.Examples(`
			pc project describe --id <project-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			projId := options.projectId
			var err error
			if projId == "" {
				projId, err = state.GetTargetProjectId()
				if err != nil {
					msg.FailMsg("No target project set and no project ID provided. Use %s to set the target project. Use %s to delete a specific project.", style.Code("pc target -p <project>"), style.Code("pc project delete -i <project-id>"))
					exit.ErrorMsg("No project ID provided, and no target project set")
				}
			}

			project, err := ac.Project.Describe(context.Background(), projId)
			if err != nil {
				msg.FailMsg("Failed to describe project %s: %s\n", projId, err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(project)
				pcio.Println(json)
			} else {
				presenters.PrintDescribeProjectTable(project)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.projectId, "id", "i", "", "ID of the project to describe")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}
