package restore

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

var (
	restoreHelp = help.Long(`
		Restore an index from a backup, and list or describe restore jobs.

		When restoring a serverless index from a backup, you can change the index name, tags,
		and deletion protection setting. All other properties of the restored index will remain
		identical to the source index, including cloud and region, dimension and similarity metric,
		and associated embedding model when restoring an index with integrated embedding.

		See: https://docs.pinecone.io/guides/manage-data/restore-an-index
	`)
)

type restoreCmdOptions struct {
	backupId           string
	name               string
	deletionProtection string
	tags               map[string]string
	json               bool
}

// RestoreJobService abstracts the Pinecone Go SDK for unit testing (runRestoreCmd)
type RestoreJobService interface {
	DescribeRestoreJob(ctx context.Context, restoreJobId string) (*pinecone.RestoreJob, error)
	ListRestoreJobs(ctx context.Context, in *pinecone.ListRestoreJobsParams) (*pinecone.RestoreJobList, error)
	CreateIndexFromBackup(ctx context.Context, in *pinecone.CreateIndexFromBackupParams) (*pinecone.CreateIndexFromBackupResponse, error)
}

func NewRestoreCmd() *cobra.Command {
	options := restoreCmdOptions{}
	cmd := &cobra.Command{
		Use:     "restore",
		Short:   "Restore an index from a backup and manage restore jobs",
		Long:    restoreHelp,
		GroupID: help.GROUP_INDEX_MANAGEMENT.ID,
		Example: help.Examples(`
			# Restore an index from a backup
			pc index restore --id backup-123 --name restored-index --tags env=prod,team=search --deletion-protection enabled

			# List restore jobs
			pc index restore list

			# Describe a restore job
			pc index restore describe --id rj-123
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runRestoreCmd(ctx, pc, options)
			if err != nil {
				msg.FailJSON(options.json, "Failed to create restore job: %s\n", err)
				exit.Error(err, "Failed to create restore job")
			}
		},
	}

	cmd.AddCommand(NewDescribeRestoreJobCmd())
	cmd.AddCommand(NewListRestoreJobsCmd())

	cmd.Flags().StringVarP(&options.backupId, "id", "i", "", "ID of the backup to restore from")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Name of the index to create from the backup")
	cmd.Flags().StringVarP(&options.deletionProtection, "deletion-protection", "d", "", "Whether to enable deletion protection on the new index (enabled|disabled)")
	cmd.Flags().StringToStringVarP(&options.tags, "tags", "t", map[string]string{}, "Tags to apply to the new index")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runRestoreCmd(ctx context.Context, svc RestoreJobService, options restoreCmdOptions) error {
	if strings.TrimSpace(options.backupId) == "" {
		return fmt.Errorf("--id is required")
	}
	if strings.TrimSpace(options.name) == "" {
		return fmt.Errorf("--name is required")
	}

	dp, err := parseDeletionProtection(options.deletionProtection)
	if err != nil {
		return err
	}

	var tags *pinecone.IndexTags
	if len(options.tags) > 0 {
		t := pinecone.IndexTags(options.tags)
		tags = &t
	}

	resp, err := svc.CreateIndexFromBackup(ctx, &pinecone.CreateIndexFromBackupParams{
		BackupId:           options.backupId,
		Name:               options.name,
		DeletionProtection: dp,
		Tags:               tags,
	})
	if err != nil {
		return err
	}

	if options.json {
		fmt.Println(text.IndentJSON(resp))
		return nil
	}

	msg.SuccessMsg("Restore job %s started for backup %s.\n", style.Emphasis(resp.RestoreJobId), style.Emphasis(options.backupId))
	msg.InfoMsg("Created index ID: %s\n", style.Emphasis(resp.IndexId))
	msg.InfoMsg("Use %s to monitor progress.\n", style.Code("pc index restore describe --id "+resp.RestoreJobId))
	return nil
}

func parseDeletionProtection(input string) (*pinecone.DeletionProtection, error) {
	if input == "" {
		return nil, nil
	}

	val := pinecone.DeletionProtection(input)
	switch val {
	case pinecone.DeletionProtectionEnabled, pinecone.DeletionProtectionDisabled:
		return &val, nil
	default:
		return nil, fmt.Errorf("invalid deletion-protection value %q, must be one of: enabled, disabled", input)
	}
}
