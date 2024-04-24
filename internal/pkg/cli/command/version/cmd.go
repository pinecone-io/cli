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
			fmt.Printf("Pinecone CLI version %s (%s)\n", build.Version, build.Date)
		},
	}

	return cmd
}
