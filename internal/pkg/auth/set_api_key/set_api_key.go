package set_api_key

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

func NewSetApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-api-key",
		Short: "Set the API key for the Pinecone CLI",
		Long: helpText,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("set-api-key called")
		},
	}

	return cmd
}