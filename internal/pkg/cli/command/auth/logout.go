package auth

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Delete all saved credentials from Pinecone CLI configuration",
		Run: func(cmd *cobra.Command, args []string) {
			config.ApiKey.Set("")
			config.SaveConfig()
			fmt.Println("Successfully logged out")
		},
	}

	return cmd
}
