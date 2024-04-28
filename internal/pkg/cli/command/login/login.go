package login

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	deviceauth "github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to Pinecone CLI",
		GroupID: help.GROUP_START.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			da := deviceauth.DeviceAuth{}
			authResponse, err := da.GetAuthResponse(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Visit %s to authenticate using code %s.\n", style.Emphasis(authResponse.VerificationURIComplete), style.HeavyEmphasis(authResponse.UserCode))
			openBrowser(authResponse.VerificationURIComplete)

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

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}
