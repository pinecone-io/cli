package pinecone

import (
	"fmt"
	"os"

	collection "github.com/pinecone-io/cli/internal/pkg/cli/command/collection"
	configCmd "github.com/pinecone-io/cli/internal/pkg/cli/command/config"
	index "github.com/pinecone-io/cli/internal/pkg/cli/command/index"
	login "github.com/pinecone-io/cli/internal/pkg/cli/command/login"
	logout "github.com/pinecone-io/cli/internal/pkg/cli/command/logout"
	org "github.com/pinecone-io/cli/internal/pkg/cli/command/org"
	project "github.com/pinecone-io/cli/internal/pkg/cli/command/project"
	target "github.com/pinecone-io/cli/internal/pkg/cli/command/target"
	version "github.com/pinecone-io/cli/internal/pkg/cli/command/version"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pinecone",
	Short: "Work seamlessly with Pinecone from the command line.",
	Long: fmt.Sprintf(`pinecone is a CLI tool for managing your Pinecone resources
	
Get started by logging in with

  %s
	`, style.CodeWithPrompt("pinecone login")),
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
	state.LoadState()

	rootCmd.SetUsageTemplate(help.HelpTemplate)

	// Getting started group
	rootCmd.AddGroup(help.GROUP_START)
	rootCmd.AddCommand(login.NewLoginCmd())
	rootCmd.AddCommand(logout.NewLogoutCmd())
	rootCmd.AddCommand(target.NewTargetCmd())

	// Management group
	rootCmd.AddGroup(help.GROUP_MANAGEMENT)
	rootCmd.AddCommand(org.NewOrgCmd())
	rootCmd.AddCommand(project.NewProjectCmd())

	// Vector database group
	rootCmd.AddGroup(help.GROUP_VECTORDB)
	rootCmd.AddCommand(index.NewIndexCmd())
	rootCmd.AddCommand(collection.NewCollectionCmd())

	// Misc group
	rootCmd.AddCommand(version.NewVersionCmd())
	rootCmd.AddCommand(configCmd.NewConfigCmd())

	// Declutter default stuff
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
}
