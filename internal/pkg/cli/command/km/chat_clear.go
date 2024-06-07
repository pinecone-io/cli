package km

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type KnowledgeModelChatClearCmdOptions struct {
	json   bool
	kmName string
}

func NewKnowledgeModelChatClearCmd() *cobra.Command {
	options := KnowledgeModelChatClearCmdOptions{}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear chat history",
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetKm.Get().Name
			if targetKm != "" {
				options.kmName = targetKm
			}
			if options.kmName == "" {
				pcio.Printf("You must target a knowledge model or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			// Reset chat history for the specified knowledge model
			chatHistory := state.ChatHist.Get()
			(*chatHistory.History)[options.kmName] = models.KnowledgeModelChat{}
			state.ChatHist.Set(&chatHistory)

			// TODO - add message for chat history reset
		},
	}

	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge model chat to clear")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
