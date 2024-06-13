package assistant

import (
	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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

			// Deleting targeted assistant, unset target
			targetAsst := state.TargetAsst.Get()
			if targetAsst.Name == options.name {
				state.TargetAsst.Clear()
				pcio.Printf("Target assistant %s deleted.\n", style.Emphasis(options.name))
				pcio.Printf("Use %s to set a new target.\n", style.Code("pinecone assistant target"))
			}

			if options.json {
				text.PrettyPrintJSON(resp)
				return
			}

			msg.SuccessMsg("Assistant %s deleted.\n", style.Emphasis(options.name))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the assistant to delete")
	cmd.MarkFlagRequired("name")
	return cmd
}
