package index

import (
	"errors"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeCmdOptions struct {
	name string
	json bool
}

func NewDescribeCmd() *cobra.Command {
	options := DescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe <name>",
		Short: "Get configuration and status information for an index",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// TODO: start interactive mode. For now just return an error.
				return errors.New("please provide an index name")
			}
			if len(args) > 1 {
				return errors.New("please provide only one index name")
			}
			if strings.TrimSpace(args[0]) == "" {
				return errors.New("index name cannot be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			options.name = args[0]
			pc := sdk.NewPineconeClient()

			idx, err := pc.DescribeIndex(cmd.Context(), options.name)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailMsg("The index %s does not exist\n", style.Emphasis(options.name))
				} else {
					msg.FailMsg("Failed to describe index %s: %s\n", style.Emphasis(options.name), err)
				}
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(idx)
				pcio.Println(json)
			} else {
				presenters.PrintDescribeIndexTable(idx)
			}
		},
	}

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
