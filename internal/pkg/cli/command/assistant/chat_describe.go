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
	json bool
	name string
}

func NewAssistantChatDescribeCmd() *cobra.Command {
	options := AssistantChatDescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe an assistant chat",
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetAsst.Get().Name
			if targetKm != "" {
				options.name = targetKm
			}
			if options.name == "" {
				pcio.Printf("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			chatHistory := state.ChatHist.Get()
			chat, ok := (*chatHistory.History)[options.name]
			if !ok {
				pcio.Printf("No chat history found for assistant %s\n", style.Emphasis(options.name))
				return
			}

			if options.json {
				text.PrettyPrintJSON(chat)
			} else {
				presenters.PrintChatHistory(chat, 100)
			}
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the assistant chat to describe")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
