package backup

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeBackupCmdOptions struct {
	id   string
	json bool
}

func NewDescribeBackupCmd() *cobra.Command {
	options := DescribeBackupCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get information on a backup",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			backup, err := pc.DescribeBackup(ctx, options.id)
			if err != nil {
				msg.FailMsg("Failed to describe backup %s: %s\n", options.id, err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(backup)
				fmt.Println(json)
			} else {
				presenters.PrintDescribeBackupTable(backup)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.id, "id", "i", "", "ID of backup to describe")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
