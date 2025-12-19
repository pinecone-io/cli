package restore

import (
	"context"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type describeRestoreJobCmdOptions struct {
	restoreJobId string
	json         bool
}

func NewDescribeRestoreJobCmd() *cobra.Command {
	options := describeRestoreJobCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a restore job by ID",
		Example: help.Examples(`
			pc pinecone backup restore describe --id rj-123
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runDescribeRestoreJobCmd(ctx, pc, options)
			if err != nil {
				msg.FailMsg("Failed to describe restore job: %s\n", err)
				exit.Error(err, "Failed to describe restore job")
			}
		},
	}

	cmd.Flags().StringVarP(&options.restoreJobId, "id", "i", "", "ID of the restore job to describe")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func runDescribeRestoreJobCmd(ctx context.Context, svc RestoreJobService, options describeRestoreJobCmdOptions) error {
	if strings.TrimSpace(options.restoreJobId) == "" {
		return pcio.Errorf("--id is required")
	}

	resp, err := svc.DescribeRestoreJob(ctx, options.restoreJobId)
	if err != nil {
		return err
	}

	if options.json {
		pcio.Println(text.IndentJSON(resp))
	} else {
		presenters.PrintRestoreJob(resp)
	}

	return nil
}
