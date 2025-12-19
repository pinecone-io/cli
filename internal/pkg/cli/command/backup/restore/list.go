package restore

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type listRestoreJobsCmdOptions struct {
	limit           int
	paginationToken string
	json            bool
}

func NewListRestoreJobsCmd() *cobra.Command {
	options := listRestoreJobsCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List restore jobs in the current project",
		Example: help.Examples(`
			pc backup restore list
			pc backup restore list --limit 5 --pagination-token token
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runListRestoreJobsCmd(ctx, pc, options)
			if err != nil {
				msg.FailMsg("Failed to list restore jobs: %s\n", err)
				exit.Error(err, "Failed to list restore jobs")
			}
		},
	}

	cmd.Flags().IntVarP(&options.limit, "limit", "l", 0, "Maximum number of restore jobs to return")
	cmd.Flags().StringVarP(&options.paginationToken, "pagination-token", "p", "", "Pagination token to continue a previous listing operation")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runListRestoreJobsCmd(ctx context.Context, svc RestoreJobService, options listRestoreJobsCmdOptions) error {
	var limit *int
	if options.limit > 0 {
		limit = &options.limit
	}

	var paginationToken *string
	if options.paginationToken != "" {
		paginationToken = &options.paginationToken
	}

	resp, err := svc.ListRestoreJobs(ctx, &pinecone.ListRestoreJobsParams{
		Limit:           limit,
		PaginationToken: paginationToken,
	})
	if err != nil {
		return err
	}

	if options.json {
		pcio.Println(text.IndentJSON(resp))
	} else {
		presenters.PrintRestoreJobList(resp)
	}

	return nil
}
