package version

import (
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/build"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type versionCmdOptions struct {
	json bool
}

func NewVersionCmd() *cobra.Command {
	options := versionCmdOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "See version information for the CLI",
		Example: help.Examples(`
			pc version
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if options.json {
				fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
					Version string `json:"version"`
					Sha     string `json:"sha"`
					Built   string `json:"built"`
				}{Version: build.Version, Sha: build.Commit, Built: build.Date}))
				return
			}

			fmt.Printf("Version: %s\n", build.Version)
			fmt.Printf("SHA: %s\n", build.Commit)
			fmt.Printf("Built: %s\n", build.Date)
		},
	}

	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}
