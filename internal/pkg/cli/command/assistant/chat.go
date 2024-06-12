package assistant

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type AssistantChatCmdOptions struct {
	name    string
	message string
	json    bool
}

func NewAssistantChatCmd() *cobra.Command {
	options := AssistantChatCmdOptions{}

	cmd := &cobra.Command{
		Use:     "chat",
		Short:   "Chat with an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetAsst.Get().Name
			if targetKm != "" {
				options.name = targetKm
			}
			if options.name == "" {
				pcio.Printf("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			// If no message is provided drop them into chat
			if options.message == "" {
				startChat(options.name)
			} else {
				// If message is provided, send it to the assistant
				sendMessage(options.name, options.message)
			}
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the assistant to chat with")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.message, "message", "m", "", "your message to the assistant")
	cmd.MarkFlagRequired("content")

	cmd.AddCommand(NewAssistantChatClearCmd())
	cmd.AddCommand(NewAssistantChatDescribeCmd())

	return cmd
}

func startChat(kmName string) {
	reader := bufio.NewReader(os.Stdin)

	// Display previous chat history up to 10 messages
	displayChatHistory(kmName, 10)

	pcio.Printf("\n\nNow chatting with assistant %s. Type your message and press Enter. Press CTRL+C to exit, or pass \"exit()\"\n\n", style.Emphasis(kmName))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		pcio.Printf("\nExiting chat with assistant %s.\n\n", style.Emphasis(kmName))
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

		checkForChatCommands(text)

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
		chatResponse, err := assistants.GetAssistantChatCompletions(kmName, message)
		if err != nil {
			exit.Error(err)
		}

		response = chatResponse

		for _, choice := range chatResponse.Choices {
			presenters.PrintAssistantChatResponse(choice.Message.Content)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func displayChatHistory(kmName string, maxNoMsgs int) {
	chatHistory := state.ChatHist.Get()
	chat, ok := (*chatHistory.History)[kmName]
	if !ok {
		pcio.Printf("No chat history found for assistant %s\n", style.Emphasis(kmName))
		return
	}

	presenters.PrintChatHistory(chat, maxNoMsgs)
}

// This function checks the input for accepted chat commands
func checkForChatCommands(text string) {
	switch text {
	case "exit()":
		pcio.Printf("Exiting chat...\n\n")
		os.Exit(0)
	default:
	}
}
