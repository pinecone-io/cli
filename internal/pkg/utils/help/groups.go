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
	GROUP_MANAGEMENT = &cobra.Group{
		ID:    "index",
		Title: style.Heading("Management Commands"),
	}
	GROUP_VECTORDB = &cobra.Group{
		ID:    "vectordb",
		Title: style.Heading("Vector Database Commands"),
	}
	GROUP_PROJECTS_API_KEYS = &cobra.Group{
		ID:    "keys",
		Title: style.Heading("API Key Management"),
	}
	GROUP_PROJECTS_CRUD = &cobra.Group{
		ID:    "projects",
		Title: style.Heading("Project Management"),
	}
	GROUP_ASSISTANT = &cobra.Group{
		ID:    "assistant",
		Title: style.Heading("Assistant Commands"),
	}
	GROUP_ASSISTANT_TARGETING = &cobra.Group{
		ID:    "assistant_targeting",
		Title: style.Heading("Target Assistant"),
	}
	GROUP_ASSISTANT_MANAGEMENT = &cobra.Group{
		ID:    "assistant_management",
		Title: style.Heading("Assistant Management"),
	}
	GROUP_ASSISTANT_OPERATIONS = &cobra.Group{
		ID:    "assistant_operations",
		Title: style.Heading("Assistant Operations"),
	}
)
