package root

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/cli/command/apiKey"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/auth"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/backup"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/collection"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/config"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/index"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/login"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/logout"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/organization"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/project"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/target"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/version"
	"github.com/pinecone-io/cli/internal/pkg/cli/command/whoami"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/spf13/cobra"
)

var (
	rootCmd        *cobra.Command
	globalOptions  GlobalOptions
	cancelRootFunc context.CancelFunc
)

type GlobalOptions struct {
	quiet   bool
	timeout time.Duration
}

func Execute() {
	//Base context: cancel on SIGINT / SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

var (
	rootHelp = help.Long(`
		The Pinecone CLI (pc) provides a complete interface for managing your
		vector database infrastructure, from authentication and project setup to
		creating indexes and managing API keys.

		GETTING STARTED

		1. Authenticate
			$ pc login
			Opens a browser to log in with your Pinecone account

		2. Set a target context (organization + project)
		   $ pc target --org "My organization" --project "My project"
		   This determines where indexes and resources will be created

		3. Create an index
		   $ pc index create --name "My index" --dimension 1536 --metric "cosine" \
		     --cloud "aws" --region "us-west-2"

		See: https://docs.pinecone.io/reference/cli/overview
	`)
)

func init() {
	// Default timeout for context.Context cancellation
	// This is applied to individual operations within subcommands through cmd.Context()
	defaultTimeout := 60 * time.Second
	globalOptions = GlobalOptions{}

	rootCmd = &cobra.Command{
		Use:   "pc",
		Short: "Manage your Pinecone vector database infrastructure from the command line",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			pcio.SetQuiet(globalOptions.quiet)

			// Apply timeout to the command context
			if globalOptions.timeout > 0 {
				ctx, cancel := context.WithTimeout(cmd.Context(), globalOptions.timeout)
				cancelRootFunc = cancel
				cmd.SetContext(ctx)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// Cancel the root context when the command completes
			if cancelRootFunc != nil {
				cancelRootFunc()
				cancelRootFunc = nil
			}
		},
		Example: help.Examples(`
		    pc login
			pc target
			pc index create --help
		`),
		Long: rootHelp,
	}

	// Help template and rendering
	rootCmd.SetHelpTemplate(help.HelpTemplate)
	help.EnableHelpRendering(rootCmd)

	// Auth group
	rootCmd.AddGroup(help.GROUP_AUTH)
	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(login.NewLoginCmd())
	rootCmd.AddCommand(logout.NewLogoutCmd())
	rootCmd.AddCommand(target.NewTargetCmd())
	rootCmd.AddCommand(whoami.NewWhoAmICmd())

	// Admin management group
	rootCmd.AddGroup(help.GROUP_ADMIN)
	rootCmd.AddCommand(organization.NewOrganizationCmd())
	rootCmd.AddCommand(project.NewProjectCmd())
	rootCmd.AddCommand(apiKey.NewAPIKeyCmd())

	// Vector database group
	rootCmd.AddGroup(help.GROUP_VECTORDB)
	rootCmd.AddCommand(index.NewIndexCmd())
	rootCmd.AddCommand(collection.NewCollectionCmd())
	rootCmd.AddCommand(backup.NewBackupCmd())

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
	rootCmd.PersistentFlags().DurationVar(&globalOptions.timeout, "timeout", defaultTimeout, "timeout for commands, defaults to 60s (0 to disable)")
}
