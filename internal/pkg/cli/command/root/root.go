package pinecone

import (
	"os"

	assistant "github.com/pinecone-io/cli/internal/pkg/cli/command/assistant"
	collection "github.com/pinecone-io/cli/internal/pkg/cli/command/collection"
	configCmd "github.com/pinecone-io/cli/internal/pkg/cli/command/config"
	index "github.com/pinecone-io/cli/internal/pkg/cli/command/index"
	login "github.com/pinecone-io/cli/internal/pkg/cli/command/login"
	logout "github.com/pinecone-io/cli/internal/pkg/cli/command/logout"
	target "github.com/pinecone-io/cli/internal/pkg/cli/command/target"
	version "github.com/pinecone-io/cli/internal/pkg/cli/command/version"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
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

func init() {
	globalOptions := GlobalOptions{}
	rootCmd = &cobra.Command{
		Use:   "pinecone",
		Short: "Work seamlessly with Pinecone from the command line.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			pcio.SetQuiet(globalOptions.quiet)
		},
		Example: help.Examples([]string{
			"pinecone login",
			"pinecone target",
			"pinecone index create-serverless --help",
		}),
		Long: pcio.Sprintf(`pinecone is a CLI tool for managing your Pinecone resources

Get started by logging in with

  %s
		`, style.CodeWithPrompt("pinecone login")),
	}

	rootCmd.SetUsageTemplate(help.HelpTemplate)

	// Getting started group
	rootCmd.AddGroup(help.GROUP_START)
	rootCmd.AddCommand(login.NewLoginCmd())
	rootCmd.AddCommand(logout.NewLogoutCmd())
	rootCmd.AddCommand(target.NewTargetCmd())
	rootCmd.AddCommand(login.NewWhoAmICmd())

	// Management group
	// rootCmd.AddGroup(help.GROUP_MANAGEMENT)
	// rootCmd.AddCommand(org.NewOrgCmd())
	// rootCmd.AddCommand(project.NewProjectCmd())

	// Vector database group
	rootCmd.AddGroup(help.GROUP_VECTORDB)
	rootCmd.AddCommand(index.NewIndexCmd())
	rootCmd.AddCommand(collection.NewCollectionCmd())

	// Assistant group
	rootCmd.AddGroup(help.GROUP_ASSISTANT)
	rootCmd.AddCommand(assistant.NewAssistantCommand())

	// Misc group
	rootCmd.AddCommand(version.NewVersionCmd())
	rootCmd.AddCommand(configCmd.NewConfigCmd())

	// Declutter default stuff
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&globalOptions.quiet, "quiet", "q", false, "suppress output")
}
