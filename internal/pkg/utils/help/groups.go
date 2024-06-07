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

	GROUP_KNOWLEDGE_MODEL = &cobra.Group{
		ID:    "km",
		Title: style.Heading("Knowledge Model Commands"),
	}

	GROUP_MANAGEMENT = &cobra.Group{
		ID:    "index",
		Title: style.Heading("Management Commands"),
	}

	GROUP_PROJECTS_API_KEYS = &cobra.Group{
		ID:    "keys",
		Title: style.Heading("API Key Management"),
	}

	GROUP_PROJECTS_CRUD = &cobra.Group{
		ID:    "projects",
		Title: style.Heading("Project Management"),
	}
	GROUP_KM_TARGETING = &cobra.Group{
		ID:    "km_targeting",
		Title: style.Heading("Target Knowledge Model"),
	}
	GROUP_KM_MANAGEMENT = &cobra.Group{
		ID:    "km_management",
		Title: style.Heading("Knowledge Model Management"),
	}
	GROUP_KM_OPERATIONS = &cobra.Group{
		ID:    "km_operations",
		Title: style.Heading("Knowledge Model Operations"),
	}
)
