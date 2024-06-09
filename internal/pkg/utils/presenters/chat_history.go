package presenters

import (
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func PrintChatHistory(chatHistory models.KnowledgeModelChat) {
	writer := NewTabWriter()
	user := os.Getenv("USER")

	pcio.Printf("%+v", chatHistory)
	messages := chatHistory.Messages
	for _, message := range messages {
		if message.Role == "user" {
			pcio.Print(style.StatusGreen(fmt.Sprintf("%s:\n", user)))
		} else {
			pcio.Printf(style.StatusYellow("Assistant:\n"))
		}
		pcio.Printf(message.Content + "\n\n")
	}

	writer.Flush()
}
