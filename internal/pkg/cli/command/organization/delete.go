package organization

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/interactive"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteOrganizationCmdOptions struct {
	organizationID   string
	skipConfirmation bool
	json             bool
}

func NewDeleteOrganizationCmd() *cobra.Command {
	options := DeleteOrganizationCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an organization by ID",
		Example: heredoc.Doc(`
		$ pc organization delete -i <organization-id>
		$ pc organization delete -i <organization-id> --skip-confirmation
		`),
		GroupID: help.GROUP_ORGANIZATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			// get the organization first
			org, err := ac.Organization.Describe(cmd.Context(), options.organizationID)
			if err != nil {
				msg.FailMsg("Failed to describe organization %s: %s\n", options.organizationID, err)
				exit.Error(err)
			}

			if !options.skipConfirmation {
				confirmDelete(org.Name, org.Id)
			}

			err = ac.Organization.Delete(cmd.Context(), options.organizationID)
			if err != nil {
				msg.FailMsg("Failed to delete organization %s: %s\n", options.organizationID, err)
				exit.Error(err)
			}

			// Clear target project if the deleted project is the target project
			if state.TargetOrg.Get().Id == options.organizationID {
				state.TargetOrg.Set(&state.TargetOrganization{
					Id:   "",
					Name: "",
				})
			}
			msg.SuccessMsg("Organization %s (ID: %s) deleted.\n", style.Emphasis(org.Name), style.Emphasis(options.organizationID))
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.organizationID, "id", "i", "", "The ID of the organization to delete")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().BoolVar(&options.skipConfirmation, "skip-confirmation", false, "Skip the deletion confirmation prompt")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func confirmDelete(organizationName string, organizationID string) {
	msg.WarnMsg("This will delete the organization %s (ID: %s).", style.Emphasis(organizationName), style.Emphasis(organizationID))
	msg.WarnMsg("This action cannot be undone.")

	question := fmt.Sprintf("Do you want to continue deleting organization '%s'?", organizationName)
	if !interactive.GetConfirmation(question) {
		msg.InfoMsg("Operation canceled.")
		exit.Success()
	}
	msg.InfoMsg("You chose to continue delete.")
}
