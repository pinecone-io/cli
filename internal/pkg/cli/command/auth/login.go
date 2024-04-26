package auth

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	deviceauth "github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Pinecone CLI",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			da := deviceauth.DeviceAuth{}
			authResponse, err := da.GetAuthResponse(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Visit %s to authenticate.\n", style.Emphasis(authResponse.VerificationURIComplete))

			style.Spinner("Waiting for authorization...", func() error {
				token, err := da.GetDeviceAccessToken(ctx, authResponse)
				if err != nil {
					return err
				}
				secrets.AccessToken.Set(token.AccessToken)
				secrets.RefreshToken.Set(token.RefreshToken)
				return nil
			})
		},
	}

	return cmd
}
