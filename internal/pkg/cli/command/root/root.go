package pinecone

import (
	"os"

	auth "github.com/pinecone-io/cli/internal/pkg/cli/command/auth"
	index "github.com/pinecone-io/cli/internal/pkg/cli/command/index"
	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pinecone",
	Short: "Work seamlessly with Pinecone from the command line.",
	Long:  `pinecone is a CLI tool managing your Pinecone resources from the command line.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.InitConfigFile()
	config.LoadConfig()

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(index.NewIndexCmd())
}
