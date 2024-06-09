package presenters

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

func PrintChatHistory(chatHistory models.KnowledgeModelChat, maxNoMsgs int) {
	writer := NewTabWriter()

	messages := truncateMessages(chatHistory.Messages, maxNoMsgs)

	for _, message := range messages {
		if message.Role == "user" {
			pcio.Print(style.StatusGreen(fmt.Sprintf("\n\n%s:\n", getUser())))
		} else {
			pcio.Printf(style.StatusYellow("\n\nAssistant:\n"))
		}
		pcio.Printf(text.WordWrapPreserveFormatting(message.Content, 80))
	}

	writer.Flush()
}

func truncateMessages(messages []models.ChatCompletionMessage, maxNoMsgs int) []models.ChatCompletionMessage {
	if maxNoMsgs <= 0 {
		maxNoMsgs = 100
	}

	if len(messages) <= maxNoMsgs {
		return messages
	}

	return messages[len(messages)-maxNoMsgs:]
}

func getUser() string {
	var user string
	if runtime.GOOS == "windows" {
		user = os.Getenv("USERNAME")
	} else {
		user = os.Getenv("USER")
	}

	if user == "" {
		user = "user"
	}
	return user
}
