package index

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type describeCmdOptions struct {
	name string
	json bool
}

func NewDescribeCmd() *cobra.Command {
	options := describeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe an index by name",
		Example: help.Examples(`
			pc index describe --name "index-name"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			idx, err := pc.DescribeIndex(cmd.Context(), options.name)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailMsg("The index %s does not exist\n", style.Emphasis(options.name))
					exit.Errorf(err, "The index %s does not exist", style.Emphasis(options.name))
				} else {
					msg.FailMsg("Failed to describe index %s: %s\n", style.Emphasis(options.name), err)
					exit.Errorf(err, "Failed to describe index %s", style.Emphasis(options.name))
				}
			}

			if options.json {
				json := text.IndentJSON(idx)
				pcio.Println(json)
			} else {
				presenters.PrintDescribeIndexTable(idx)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe")
	_ = cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	return cmd
}
