package project

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type createProjectCmdOptions struct {
	name                    string
	forceEncryptionWithCMEK bool
	maxPods                 int
	target                  bool
	json                    bool
}

var (
	createHelp = help.Long(`
		Create a new project in your target organization.
		
		In Pinecone, projects are organizational containers where you create and
		manage indexes. All indexes must belong to a project, and each project has
		its own API keys.
		
		To target the newly created project, include the '--target' flag.
		See: https://docs.pinecone.io/guides/projects/manage-projects
	`)
)

func NewCreateProjectCmd() *cobra.Command {
	options := createProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a project for the target organization determined by user credentials",
		Long:    createHelp,
		GroupID: help.GROUP_PROJECTS.ID,
		Example: help.Examples(`
			pc project create --name "demo-project" --max-pods 10 --force-encryption
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ac := sdk.NewPineconeAdminClient(ctx)

			createParams := &pinecone.CreateProjectParams{}
			if options.name != "" {
				createParams.Name = options.name
			}
			if options.maxPods > 0 {
				createParams.MaxPods = &options.maxPods
			}
			if options.forceEncryptionWithCMEK {
				createParams.ForceEncryptionWithCmek = &options.forceEncryptionWithCMEK
			}

			proj, err := ac.Project.Create(ctx, createParams)
			if err != nil {
				msg.FailMsg("Failed to create project %s: %s\n", style.Emphasis(options.name), err)
				exit.Errorf(err, "Failed to create project %s", style.Emphasis(options.name))
			}

			if options.json {
				json := text.IndentJSON(proj)
				pcio.Println(json)
				return
			}

			msg.SuccessMsg("Project %s created successfully.\n", style.Emphasis(proj.Name))
			presenters.PrintDescribeProjectTable(proj)

			// If the user has requested to swap targeting to the newly created project
			if options.target {
				state.TargetProj.Set(state.TargetProject{
					Name: proj.Name,
					Id:   proj.Id,
				})
				msg.SuccessMsg("Target project set to %s", style.Emphasis(proj.Name))
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Name of the project")
	_ = cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().IntVarP(&options.maxPods, "max-pods", "p", 5, "Maximum number of Pods that can be created in the project across all indexes")
	cmd.Flags().BoolVar(&options.forceEncryptionWithCMEK, "force-encryption", false, "Whether to force encryption with a customer-managed encryption key (CMEK). Default is 'false'. Once enabled, CMEK encryption cannot be disabled.")
	cmd.Flags().BoolVar(&options.target, "target", false, "Target the newly created project")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")
	return cmd
}
