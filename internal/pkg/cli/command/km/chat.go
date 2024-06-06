package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type KnowledgeModelChatCmdOptions struct {
	kmName  string
	content string
	json    bool
}

func NewKnowledgeModelChatCmd() *cobra.Command {
	options := KnowledgeModelChatCmdOptions{}

	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Chat with a knowledge model",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := knowledge.GetKnowledgeModelSearchCompletions(options.kmName, options.content)
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
		},
	}
	// required flags
	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge base to describe")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&options.content, "content", "c", "", "your message to the knowledge model")
	cmd.MarkFlagRequired("content")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
