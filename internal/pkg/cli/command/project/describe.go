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
		Short:   "Describe a project by ID",
		GroupID: help.GROUP_PROJECTS.ID,
		Example: heredoc.Doc(`
		$ pc project describe -i <project-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			project, err := ac.Project.Describe(context.Background(), options.projectId)
			if err != nil {
				msg.FailMsg("Failed to describe project %s: %s\n", options.projectId, err)
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
