package organization

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
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type deleteOrganizationCmdOptions struct {
	organizationID   string
	skipConfirmation bool
	json             bool
}

type deleteOrganizationService interface {
	Delete(ctx context.Context, id string) error
}

func NewDeleteOrganizationCmd() *cobra.Command {
	options := deleteOrganizationCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an organization by ID",
		Example: help.Examples(`
			pc organization delete --id "organization-id""
			pc organization delete --id "organization-id" --skip-confirmation
		`),
		GroupID: help.GROUP_ORGANIZATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ac := sdk.NewPineconeAdminClient(ctx)

			// get the organization first
			org, err := ac.Organization.Describe(cmd.Context(), options.organizationID)
			if err != nil {
				msg.FailMsg("Failed to describe organization %s: %s\n", options.organizationID, err)
				exit.Errorf(err, "Failed to describe organization %s", style.Emphasis(options.organizationID))
			}

			if !options.skipConfirmation {
				confirmDelete(org.Name, org.Id)
			}

			err = runDeleteOrganizationCmd(ctx, ac.Organization, options, org.Name, org.Id)
			if err != nil {
				msg.FailMsg("Failed to delete organization %s: %s\n", options.organizationID, err)
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
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runDeleteOrganizationCmd(ctx context.Context, svc deleteOrganizationService, opts deleteOrganizationCmdOptions, name, id string) error {
	if err := svc.Delete(ctx, id); err != nil {
		return err
	}

	if opts.json {
		fmt.Println(text.IndentJSON(struct {
			Deleted bool   `json:"deleted"`
			Name    string `json:"name"`
			Id      string `json:"id"`
		}{Deleted: true, Name: name, Id: id}))
		return nil
	}

	msg.SuccessMsg("Organization %s (ID: %s) deleted.\n", style.Emphasis(name), style.Emphasis(id))
	return nil
}

func confirmDelete(organizationName string, organizationID string) {
	msg.WarnMsg("This will delete the organization %s (ID: %s).", style.Emphasis(organizationName), style.Emphasis(organizationID))
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

	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" {
		msg.InfoMsg("You chose to continue delete.")
	} else {
		msg.InfoMsg("Operation canceled.")
		exit.Success()
	}
}
