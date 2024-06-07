package km

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type KnowledgeModelChatDescribeCmdOptions struct {
	json   bool
	kmName string
}

func NewKnowledgeModelChatDescribeCmd() *cobra.Command {
	options := KnowledgeModelChatDescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a knowledge model chat",
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetKm.Get().Name
			if targetKm != "" {
				options.kmName = targetKm
			}
			if options.kmName == "" {
				pcio.Printf("You must target a knowledge model or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			chatHistory := state.ChatHist.Get()
			chat, ok := (*chatHistory.History)[options.kmName]

			if !ok {
				pcio.Printf("No chat history found for knowledge model %s\n", style.Emphasis(options.kmName))
				return
			}

			if options.json {
				text.PrettyPrintJSON(chat)
			} else {
				presenters.PrintChatHistory(chat)
			}
		},
	}

	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge base to describe")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
