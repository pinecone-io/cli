package presenters

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func PrintKnowledgeChatResponse(response string) {
	writer := NewTabWriter()

	pcio.Printf(response)
	// columns := []string{"RESPONSE", "SCORE"}
	// header := strings.Join(columns, "\t") + "\n"
	// pcio.Fprint(writer, header)

	// for _, r := range resp.Responses {
	// 	pcio.Fprintf(writer, "%s\t%f\n", r.Response, r.Score)
	// }

	writer.Flush()
}
