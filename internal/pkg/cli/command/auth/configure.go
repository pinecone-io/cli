package auth

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type InitCmdOptions struct {
	clientID            string
	clientSecret        string
	readSecretFromStdin bool
	promptIfMissing     bool
}

func NewConfigureCmd() *cobra.Command {
	options := InitCmdOptions{}

	cmd := &cobra.Command{
		Use:     "configure",
		Short:   "Initilize authentication credentials for the Pinecone CLI",
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
				out = io.Discard
			}

			Run(cmd.Context(), IO{
				In:  cmd.InOrStdin(),
				Out: out,
				Err: cmd.ErrOrStderr(),
			}, options)
		},
	}

	cmd.Flags().StringVar(&options.clientID, "client-id", "", "client id for the Pinecone CLI")
	cmd.Flags().StringVar(&options.clientSecret, "client-secret", "", "client secret for the Pinecone CLI")
	cmd.Flags().BoolVar(&options.readSecretFromStdin, "client-secret-stdin", false, "read client secret from stdin")
	cmd.Flags().BoolVar(&options.promptIfMissing, "prompt-if-missing", false, "prompt for missing credentials if not provided")

	return cmd
}

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func Run(ctx context.Context, io IO, opts InitCmdOptions) {
	secret := strings.TrimSpace(opts.clientSecret)
	if secret == "" {
		if opts.readSecretFromStdin {
			secretBytes, err := ioReadAll(io.In)
			if err != nil {
				log.Error().Err(err).Msg("Error reading client secret from stdin")
				exit.Error(pcio.Errorf("error reading client secret from stdin: %w", err))
			}
			secret = string(secretBytes)
		} else if opts.promptIfMissing && isTerminal(os.Stdin) {
			pcio.Fprint(io.Out, "Client Secret: ")
			secretBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Error().Err(err).Msg("Error reading client secret from terminal")
				exit.Error(pcio.Errorf("error reading client secret from terminal: %w", err))
			}
			secret = string(secretBytes)
		}
	}

	if secret == "" {
		log.Error().Msg("Error configuring authentication credentials")
		msg.FailMsg("Client secret is required (use %s or %s to provide it)", style.Emphasis("--client-secret"), style.Emphasis("--client-secret-stdin"))
		exit.Error(pcio.Errorf("client secret is required"))
		return
	}

	// store values
	secrets.ClientId.Set(strings.TrimSpace(opts.clientID))
	secrets.ClientSecret.Set(secret)
}

func ioReadAll(r io.Reader) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	var buf strings.Builder
	tmp := make([]byte, 4096)
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return []byte(buf.String()), nil
}

func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
