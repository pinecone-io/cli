package auth

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/build"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Short:  "See version information for the CLI",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", build.Version)
			fmt.Printf("SHA: %s\n", build.Commit)
			fmt.Printf("Built: %s\n", build.Date)
		},
	}

	return cmd
}
