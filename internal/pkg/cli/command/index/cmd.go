package index

import (
	"github.com/spf13/cobra"
)

var helpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

func NewIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "index <command>",
		Short:   "Work with indexes",
		Long: helpText,
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.AddCommand(NewDescribeCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateServerlessCmd())

	return cmd
}

