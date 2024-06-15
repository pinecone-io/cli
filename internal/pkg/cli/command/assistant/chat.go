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
	stream  bool
	json    bool
}

func NewAssistantChatCmd() *cobra.Command {
	options := AssistantChatCmdOptions{}

	cmd := &cobra.Command{
		Use:     "chat",
		Short:   "Chat with an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetAsst := state.TargetAsst.Get().Name
			if targetAsst != "" {
				options.name = targetAsst
			}
			if options.name == "" {
				pcio.Printf("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			// If no message is provided drop them into interactive chat
			if options.message == "" {
				startInteractiveChat(options.name)
			} else {
				// If message is provided, send it to the assistant
				sendMessage(options.name, options.message, options.stream)
			}
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the assistant to chat with")
	cmd.Flags().StringVarP(&options.message, "message", "m", "", "your message to the assistant")
	cmd.Flags().BoolVarP(&options.stream, "stream", "s", false, "stream chat message responses")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.MarkFlagRequired("content")

	cmd.AddCommand(NewAssistantChatClearCmd())
	cmd.AddCommand(NewAssistantChatDescribeCmd())

	return cmd
}

func startInteractiveChat(asstName string) {
	reader := bufio.NewReader(os.Stdin)

	// Display previous chat history up to 10 messages
	displayChatHistory(asstName, 10)

	pcio.Printf("\n\nNow chatting with assistant %s. Type your message and press Enter. Press CTRL+C to exit, or pass \"exit()\"\n\n", style.Emphasis(asstName))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		pcio.Printf("\nExiting chat with assistant %s.\n\n", style.Emphasis(asstName))
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
			// Stream here since we're in interactive chat mode
			_, err := sendMessage(asstName, text, true)
			if err != nil {
				pcio.Printf("Error sending message: %s\n", err)
				continue
			}
		}
	}
}

func sendMessage(asstName string, message string, stream bool) (*models.ChatCompletionModel, error) {
	var chatGetter func() (*models.ChatCompletionModel, error)
	if stream {
		chatGetter = streamChatResponse(asstName, message, stream)
	} else {
		chatGetter = getChatResponse(asstName, message, stream)
	}

	chatResp, err := chatGetter()
	if err != nil {
		return nil, err
	}

	return chatResp, nil
}

func getChatResponse(asstName string, message string, stream bool) func() (*models.ChatCompletionModel, error) {
	return func() (*models.ChatCompletionModel, error) {
		chatResponse, err := assistants.GetAssistantChatCompletions(asstName, message, stream)
		if err != nil {
			exit.Error(err)
		}

		for _, choice := range chatResponse.Choices {
			presenters.PrintAssistantChatResponse(choice.Message.Content)
		}
		return chatResponse, nil
	}
}

func streamChatResponse(asstName string, message string, stream bool) func() (*models.ChatCompletionModel, error) {
	return func() (*models.ChatCompletionModel, error) {
		chatResponse, err := assistants.GetAssistantChatCompletions(asstName, message, stream)
		if err != nil {
			exit.Error(err)
		}
		// We don't print the chat response since it's printed while streamed in assistants.PostAndStreamChatResponse

		return chatResponse, nil
	}
}

func displayChatHistory(asstName string, maxNoMsgs int) {
	chatHistory := state.ChatHist.Get()
	chat, ok := (*chatHistory.History)[asstName]
	if !ok {
		pcio.Printf("No chat history found for assistant %s\n", style.Emphasis(asstName))
		return
	}

	presenters.PrintChatHistory(chat, maxNoMsgs)
}

// Checks the input for accepted chat commands
func checkForChatCommands(text string) {
	switch text {
	case "exit()":
		pcio.Printf("Exiting chat...\n\n")
		os.Exit(0)
	default:
	}
}
