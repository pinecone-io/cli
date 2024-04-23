package presenters

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func ColorizeBool(b bool) string {
	if b {
		return style.StatusGreen("true")
	}
	return style.StatusRed("false")
}
