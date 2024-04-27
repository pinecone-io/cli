package help

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func Examples(examples []string) string {
	const pad = "  "
	for i := range examples {
		examples[i] = pad + style.CodeWithPrompt(examples[i])
	}
	return strings.Join(examples, "\n")
}
