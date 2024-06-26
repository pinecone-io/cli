package assistant

import (
	assistant "github.com/pinecone-io/cli/internal/pkg/cli/command/assistant/chat"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

var helpText = text.WordWrap(`Pinecone Assistant Engine is a context engine to store and retrieve relevant knowledge 
    from millions of documents at scale. This API supports creating and managing assistants.`, 80)

func NewAssistantCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assistant <command>",
		Short:   "Work with assistants",
		Long:    helpText,
		GroupID: help.GROUP_ASSISTANT.ID,
	}

	// Targeting
	cmd.AddGroup(help.GROUP_ASSISTANT_TARGETING)
	cmd.AddCommand(NewAssistantTargetCmd())

	// Assistant Management
	cmd.AddGroup(help.GROUP_ASSISTANT_MANAGEMENT)
	cmd.AddCommand(NewCreateAssistantCmd())
	cmd.AddCommand(NewListAssistantsCmd())
	cmd.AddCommand(NewDeleteAssistantCmd())

	// Assistant Operations
	cmd.AddGroup(help.GROUP_ASSISTANT_OPERATIONS)
	cmd.AddCommand(NewDescribeAssistantCmd())
	cmd.AddCommand(assistant.NewAssistantChatCmd())
	cmd.AddCommand(NewListAssistantFilesCmd())
	cmd.AddCommand(NewDeleteAssistantFileCmd())
	cmd.AddCommand(NewUploadAssistantFileCmd())
	cmd.AddCommand(NewDescribeAssistantFileCmd())

	return cmd
}
