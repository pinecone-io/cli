package root

import (
	"os"

	"github.com/pinecone-io/cli/internal/pkg/cli/command/apiKey"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/auth"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/collection"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/config"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/index"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/login"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/logout"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/organization"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/project"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/target"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/version"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

type GlobalOptions struct {
	quiet bool
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	globalOptions := GlobalOptions{}
	rootCmd = &cobra.Command{
		Use:   "pc",
		Short: "Work seamlessly with Pinecone from the command line.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			pcio.SetQuiet(globalOptions.quiet)
		},
		Example: help.Examples(`
		    pc login
			pc target
			pc index create --help
		`),
		Long: help.Long(`
			pc is a CLI tool for managing your Pinecone resources

			Get started by logging in with pc login
		`),
	}

	rootCmd.SetHelpTemplate(help.HelpTemplate)
	help.EnableHelpRendering(rootCmd)

	// Auth group
	rootCmd.AddGroup(help.GROUP_AUTH)
	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(login.NewLoginCmd())
	rootCmd.AddCommand(logout.NewLogoutCmd())
	rootCmd.AddCommand(target.NewTargetCmd())
	rootCmd.AddCommand(login.NewWhoAmICmd())

	// Admin management group
	rootCmd.AddGroup(help.GROUP_ADMIN)
	rootCmd.AddCommand(organization.NewOrganizationCmd())
	rootCmd.AddCommand(project.NewProjectCmd())
	rootCmd.AddCommand(apiKey.NewAPIKeyCmd())

	// Vector database group
	rootCmd.AddGroup(help.GROUP_VECTORDB)
	rootCmd.AddCommand(index.NewIndexCmd())
	rootCmd.AddCommand(collection.NewCollectionCmd())

	// Misc group
	rootCmd.AddCommand(version.NewVersionCmd())
	rootCmd.AddCommand(config.NewConfigCmd())

	// Declutter default stuff
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&globalOptions.quiet, "quiet", "q", false, "suppress output")
}
