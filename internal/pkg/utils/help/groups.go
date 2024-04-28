package help

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	cobra "github.com/spf13/cobra"
)

var (
	GROUP_START = &cobra.Group{
		ID:    "getting-started",
		Title: style.Heading("Auth Commands"),
	}

	GROUP_VECTORDB = &cobra.Group{
		ID:    "vectordb",
		Title: style.Heading("Vector Database Commands"),
	}

	GROUP_MANAGEMENT = &cobra.Group{
		ID:    "index",
		Title: style.Heading("Management Commands"),
	}
)