package pinecone

import (
	"os"

	"github.com/spf13/cobra"
	auth "github.com/pinecone-io/cli/internal/pkg/auth"
	index "github.com/pinecone-io/cli/internal/pkg/index"
)

var rootCmd = &cobra.Command{
	Use:   "pinecone",
	Short: "Work seamlessly with Pinecone from the command line.",
	Long: `pinecone is a CLI tool managing your Pinecone resources from the command line.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(index.NewIndexCmd())
}


