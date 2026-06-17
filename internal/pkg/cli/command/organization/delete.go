package organization

import (
	"context"
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/confirm"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type deleteOrganizationCmdOptions struct {
	organizationID   string
	skipConfirmation bool
	json             bool
}

// DeleteOrganizationService abstracts the Pinecone Go SDK for unit testing (runDeleteOrganizationCmd)
type DeleteOrganizationService interface {
	Delete(ctx context.Context, id string) error
}

func NewDeleteOrganizationCmd() *cobra.Command {
	options := deleteOrganizationCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an organization by ID",
		Example: help.Examples(`
			pc organization delete --id "organization-id"
			pc organization delete --id "organization-id" --skip-confirmation
		`),
		GroupID: help.GROUP_ORGANIZATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ac := sdk.NewPineconeAdminClient(ctx)

			// get the organization first
			org, err := ac.Organization.Describe(cmd.Context(), options.organizationID)
			if err != nil {
				msg.FailJSON(options.json, "Failed to describe organization %s: %s\n", options.organizationID, err)
				exit.Errorf(err, "Failed to describe organization %s", style.Emphasis(options.organizationID))
			}

			if !options.skipConfirmation && !options.json {
				confirm.Deletion(
					fmt.Sprintf("This will delete the organization %s (ID: %s).", style.Emphasis(org.Name), style.Emphasis(org.Id)),
					"This action cannot be undone.",
				)
			}

			err = runDeleteOrganizationCmd(ctx, ac.Organization, options, org.Name, org.Id)
			if err != nil {
				msg.FailJSON(options.json, "Failed to delete organization %s: %s\n", options.organizationID, err)
				exit.Errorf(err, "Failed to delete organization %s", style.Emphasis(options.organizationID))
			}

			// Clear target org if the deleted org is the target org
			if state.TargetOrg.Get().Id == options.organizationID {
				state.TargetOrg.Set(state.TargetOrganization{
					Id:   "",
					Name: "",
				})
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.organizationID, "id", "i", "", "The ID of the organization to delete")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().BoolVar(&options.skipConfirmation, "skip-confirmation", false, "Skip the deletion confirmation prompt")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON (also skips confirmation prompt)")

	return cmd
}

func runDeleteOrganizationCmd(ctx context.Context, svc DeleteOrganizationService, opts deleteOrganizationCmdOptions, name, id string) error {
	if err := svc.Delete(ctx, id); err != nil {
		return err
	}

	if opts.json {
		fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
			Deleted bool   `json:"deleted"`
			Name    string `json:"name"`
			Id      string `json:"id"`
		}{Deleted: true, Name: name, Id: id}))
		return nil
	}

	msg.SuccessMsg("Organization %s (ID: %s) deleted.\n", style.Emphasis(name), style.Emphasis(id))
	return nil
}
