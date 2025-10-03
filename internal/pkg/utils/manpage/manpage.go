package manpage

import (
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func GenerateManPages(rootCmd *cobra.Command, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Disable ANSI colors in man pages
	prevColor := config.Color.Get()
	config.Color.Set(false)
	defer config.Color.Set(prevColor)

	// Configure header
	header := &doc.GenManHeader{
		Title:   "PC",
		Section: "1",
		Source:  "Pinecone CLI",
		Manual:  "Pinecone CLI Manual",
	}

	// Generate man pages for root command and all subcommands
	return doc.GenManTree(rootCmd, header, outputDir)
}
