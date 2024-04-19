package logout

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Delete all saved credentials from Pinecone CLI configuration",
		Long: helpText,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("logout called")
		},
	}

	return cmd
}