package project

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type CreateProjectCmdOptions struct {
	organizationID          string
	name                    string
	forceEncryptionWithCMEK bool
	maxPods                 int
	json                    bool
}

func NewCreateProjectCmd() *cobra.Command {
	options := CreateProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a project for a specific organization by ID or the target organization",
		GroupID: help.GROUP_PROJECTS.ID,
		Example: heredoc.Doc(`
		$ pc target -o "my-organization-name"
		$ pc project create --name "demo-project" --max-pods 10 --force-encryption
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			_, err := state.GetTargetOrgId()
			if err != nil {
				msg.FailMsg("No target organization set. Use %s to set the organization context.", style.Code("pc target -o <org>"))
				cmd.Help()
				exit.ErrorMsg("No organization context set")
			}

			proj, err := ac.Project.Create(context.Background(), &pinecone.CreateProjectParams{
				Name:                    options.name,
				MaxPods:                 &options.maxPods,
				ForceEncryptionWithCmek: &options.forceEncryptionWithCMEK,
			})
			if err != nil {
				msg.FailMsg("Failed to create project %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(proj)
				pcio.Println(json)
				return
			}

			msg.SuccessMsg("Project %s created successfully.\n", style.Emphasis(proj.Name))
			presenters.PrintDescribeProjectTable(proj)
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Name of the project")
	_ = cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().StringVarP(&options.organizationID, "id", "i", "", "The ID of the organization to create the project in if not the target organization")
	cmd.Flags().IntVarP(&options.maxPods, "max-pods", "p", 5, "Maximum number of Pods that can be created in the project across all indexes")
	cmd.Flags().BoolVar(&options.forceEncryptionWithCMEK, "force-encryption", false, "Whether to force encryption with a customer-managed encryption key (CMEK). Default is 'false'. Once enabled, CMEK encryption cannot be disabled.")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")
	return cmd
}
