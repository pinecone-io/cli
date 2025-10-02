package version

import (
	"github.com/pinecone-io/cli/internal/build"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "See version information for the CLI",
		Example: help.Examples(`
			pc version
		`),
		Run: func(cmd *cobra.Command, args []string) {
			pcio.Printf("Version: %s\n", build.Version)
			pcio.Printf("SHA: %s\n", build.Commit)
			pcio.Printf("Built: %s\n", build.Date)
		},
	}

	return cmd
}
