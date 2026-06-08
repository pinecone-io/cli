package index

import (
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type describeCmdOptions struct {
	indexName string
	json      bool
}

func NewDescribeCmd() *cobra.Command {
	options := describeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe an index by name",
		Example: help.Examples(`
			pc index describe --index-name "index-name"
		`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("index-name") && !cmd.Flags().Changed("name") {
				return fmt.Errorf("required flag(s) \"index-name\" not set")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			idx, err := pc.DescribeIndex(cmd.Context(), options.indexName)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailJSON(options.json, "The index %s does not exist\n", style.Emphasis(options.indexName))
					exit.Errorf(err, "The index %s does not exist", style.Emphasis(options.indexName))
				} else {
					msg.FailJSON(options.json, "Failed to describe index %s: %s\n", style.Emphasis(options.indexName), err)
					exit.Errorf(err, "Failed to describe index %s", style.Emphasis(options.indexName))
				}
			}

			if options.json {
				json := text.IndentJSON(idx)
				fmt.Fprintln(os.Stdout, json)
			} else {
				presenters.PrintDescribeIndexTable(idx)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.indexName, "index-name", "i", "", "name of index to describe")
	cmd.Flags().StringVarP(&options.indexName, "name", "n", "", "name of index to describe")
	_ = cmd.Flags().MarkDeprecated("name", "use --index-name instead")

	// optional flags
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	return cmd
}
