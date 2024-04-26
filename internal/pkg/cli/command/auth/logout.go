package auth

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Delete all saved credentials from Pinecone CLI configuration",
		Run: func(cmd *cobra.Command, args []string) {
			secrets.ApiKey.Set("")
			secrets.AccessToken.Set("")
			secrets.SaveSecrets()
			fmt.Println("✅ Secrets cleared.")
		},
	}

	return cmd
}
