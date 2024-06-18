package assistant

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type AssistantChatDescribeCmdOptions struct {
	assistant string
	id        bool
	json      bool
}

func NewAssistantChatDescribeCmd() *cobra.Command {
	options := AssistantChatDescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe an assistant chat",
		Run: func(cmd *cobra.Command, args []string) {
			targetAsst := state.TargetAsst.Get().Name
			if targetAsst != "" {
				options.assistant = targetAsst
			}
			if options.assistant == "" {
				pcio.Printf("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			chatHistory := state.ChatHist.Get()
			chat, ok := (*chatHistory.History)[options.assistant]
			if !ok {
				pcio.Printf("No chat history found for assistant %s\n", style.Emphasis(options.assistant))
				return
			}

			// If the chat ID was requested print that
			if options.id {
				pcio.Printf("id: %s\n", chat.Id)
				return
			}

			// Otherwise print the chat history
			if options.json {
				text.PrettyPrintJSON(chat)
			} else {
				presenters.PrintChatHistory(chat, 100)
			}
		},
	}

	cmd.Flags().StringVarP(&options.assistant, "assistant", "a", "", "name of the assistant chat to describe")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().BoolVar(&options.id, "id", false, "output the ID of the chat")

	return cmd
}
