package presenters

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func PrintKnowledgeChatResponse(response string) {
	writer := NewTabWriter()

	pcio.Printf(response)
	// TODO - implement better display UX for chat responses

	writer.Flush()
}
