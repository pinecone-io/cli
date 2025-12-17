package testutils

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

// Silences pcio output for tests and, and returns a function to restore it
func SilenceOutput() func() {
	pcio.SetQuiet(true)
	return func() { pcio.SetQuiet(false) }
}
