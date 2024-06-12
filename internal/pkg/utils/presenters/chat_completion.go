package presenters

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

func PrintAssistantChatResponse(response string) {
	writer := NewTabWriter()

	pcio.Fprintf(writer, style.StatusYellow("\n\nAssistant:\n"))
	pcio.Fprintf(writer, text.WordWrapPreserveFormatting(response, 80))

	writer.Flush()
}
