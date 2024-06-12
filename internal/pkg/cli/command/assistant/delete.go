package assistant

import (
	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DeleteAssistantCmdOptions struct {
	name string
	json bool
}

func NewDeleteAssistantCmd() *cobra.Command {
	options := DeleteAssistantCmdOptions{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an assistant",
		GroupID: help.GROUP_ASSISTANT_MANAGEMENT.ID,
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := assistants.DeleteAssistant(options.name)
			if err != nil {
				msg.FailMsg("Failed to delete assistant %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(resp)
				return
			}

			msg.SuccessMsg("Assistant %s deleted.\n", style.Emphasis(options.name))

			// TODO - check to see if the current target is the delete assistant, if so we need to clear it
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the assistant to delete")
	cmd.MarkFlagRequired("name")
	return cmd
}
