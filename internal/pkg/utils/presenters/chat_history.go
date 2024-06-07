package presenters

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func PrintChatHistory(chatHistory models.KnowledgeModelChat) {
	writer := NewTabWriter()

	pcio.Printf("%+v", chatHistory)
	// TODO - implement better display UX for chat history

	writer.Flush()
}
