package help

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	cobra "github.com/spf13/cobra"
)

var (
	GROUP_AUTH = &cobra.Group{
		ID:    "auth",
		Title: style.Heading("Authentication Commands"),
	}
	GROUP_ADMIN = &cobra.Group{
		ID:    "admin",
		Title: style.Heading("Admin Management Commands"),
	}
	GROUP_API_KEYS = &cobra.Group{
		ID:    "api-keys",
		Title: style.Heading("API Key Management Commands"),
	}
	GROUP_PROJECTS = &cobra.Group{
		ID:    "projects",
		Title: style.Heading("Project Management Commands"),
	}
	GROUP_ORGANIZATIONS = &cobra.Group{
		ID:    "organizations",
		Title: style.Heading("Organization Management Commands "),
	}
	GROUP_VECTORDB = &cobra.Group{
		ID:    "vectordb",
		Title: style.Heading("Vector Database Commands"),
	}
	GROUP_INDEX_DATA = &cobra.Group{
		ID:    "index-data",
		Title: style.Heading("Index Data Commands"),
	}
)
