package pinecone

import (
	"os"

	auth "github.com/pinecone-io/cli/internal/pkg/cli/command/auth"
	collection "github.com/pinecone-io/cli/internal/pkg/cli/command/collection"
	configCmd "github.com/pinecone-io/cli/internal/pkg/cli/command/config"
	index "github.com/pinecone-io/cli/internal/pkg/cli/command/index"
	org "github.com/pinecone-io/cli/internal/pkg/cli/command/org"
	project "github.com/pinecone-io/cli/internal/pkg/cli/command/project"
	version "github.com/pinecone-io/cli/internal/pkg/cli/command/version"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pinecone",
	Short: "Work seamlessly with Pinecone from the command line.",
	Long:  `pinecone is a CLI tool for managing your Pinecone resources`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.LoadConfig()
	secrets.LoadSecrets()

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(index.NewIndexCmd())
	rootCmd.AddCommand(collection.NewCollectionCmd())
	rootCmd.AddCommand(configCmd.NewConfigCmd())
	rootCmd.AddCommand(version.NewVersionCmd())
	rootCmd.AddCommand(org.NewOrgCmd())
	rootCmd.AddCommand(project.NewProjectCmd())
}
