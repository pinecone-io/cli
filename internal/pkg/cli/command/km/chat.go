package km

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
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

			// If no message is provided drop them into chat
			if options.message == "" {
				startChat(options.kmName)
			} else {
				// If message is provided, send it to the knowledge model
				sendMessage(options.kmName, options.message)
			}
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

func startChat(kmName string) {
	reader := bufio.NewReader(os.Stdin)
	pcio.Printf("Now chatting with knowledge model %s. Type your message and press Enter. Press CTRL+C to exit.\n\n", style.Emphasis(kmName))

	// Display previous chat history
	displayChatHistory(kmName)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		pcio.Printf("\nExiting chat with knowledge model %s.\n\n", style.Emphasis(kmName))
		os.Exit(0)
	}()

	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			pcio.Printf("Error reading input: %s\n", err)
			continue
		}

		text = strings.TrimSpace(text)

		if text != "" {
			_, err := sendMessage(kmName, text)
			if err != nil {
				pcio.Printf("Error sending message: %s\n", err)
				continue
			}
		}
	}
}

func sendMessage(kmName string, message string) (*models.ChatCompletionModel, error) {
	response := &models.ChatCompletionModel{}

	err := style.Spinner("", func() error {
		chatResponse, err := knowledge.GetKnowledgeModelSearchCompletions(kmName, message)
		if err != nil {
			exit.Error(err)
		}

		response = chatResponse

		for _, choice := range chatResponse.Choices {
			presenters.PrintKnowledgeChatResponse(choice.Message.Content)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func displayChatHistory(kmName string) {
	chatHistory := state.ChatHist.Get()
	chat, ok := (*chatHistory.History)[kmName]
	if !ok {
		pcio.Printf("No chat history found for knowledge model %s\n", style.Emphasis(kmName))
		return
	}

	presenters.PrintChatHistory(chat)
}
