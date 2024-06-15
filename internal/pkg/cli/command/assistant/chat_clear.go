package assistant

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type AssistantChatClearCmdOptions struct {
	assistant string
	json      bool
}

func NewAssistantChatClearCmd() *cobra.Command {
	options := AssistantChatClearCmdOptions{}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear chat history",
		Run: func(cmd *cobra.Command, args []string) {
			targetAsst := state.TargetAsst.Get().Name
			if targetAsst != "" {
				options.assistant = targetAsst
			}
			if options.assistant == "" {
				pcio.Printf("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			// Reset chat history for the specified assistant
			chatHistory := state.ChatHist.Get()
			(*chatHistory.History)[options.assistant] = models.AssistantChat{}
			state.ChatHist.Set(&chatHistory)

			pcio.Printf("Chat history for assistant %s cleared.\n", options.assistant)
		},
	}

	cmd.Flags().StringVarP(&options.assistant, "assistant", "a", "", "name of the assistant chat to clear")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
