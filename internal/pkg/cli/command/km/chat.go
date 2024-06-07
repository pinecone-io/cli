package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type KnowledgeModelChatCmdOptions struct {
	kmName  string
	message string
	json    bool
}

func NewKnowledgeModelChatCmd() *cobra.Command {
	options := KnowledgeModelChatCmdOptions{}

	cmd := &cobra.Command{
		Use:     "chat",
		Short:   "Chat with a knowledge model",
		GroupID: help.GROUP_KM_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetKm.Get().Name
			if targetKm != "" {
				options.kmName = targetKm
			}
			if options.kmName == "" {
				pcio.Printf("You must target a knowledge model or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			style.Spinner("", func() error {
				resp, err := knowledge.GetKnowledgeModelSearchCompletions(options.kmName, options.message)
				if err != nil {
					exit.Error(err)
				}

				if options.json {
					text.PrettyPrintJSON(resp)
				} else {
					for _, choice := range resp.Choices {
						presenters.PrintKnowledgeChatResponse(choice.Message.Content)
					}
				}
				return nil
			})
		},
	}

	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge model to chat with")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.message, "message", "m", "", "your message to the knowledge model")
	cmd.MarkFlagRequired("content")

	cmd.AddCommand(NewKnowledgeModelChatClearCmd())
	cmd.AddCommand(NewKnowledgeModelChatDescribeCmd())

	return cmd
}
