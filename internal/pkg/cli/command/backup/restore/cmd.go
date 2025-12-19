package restore

import (
	"context"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

var (
	restoreJobHelp = help.Long(`
		Inspect restore jobs that were initiated when creating an index from a backup.
	`)
)

type restoreJobCmdOptions struct {
	backupId           string
	name               string
	deletionProtection string
	tags               map[string]string
	json               bool
}

type RestoreJobService interface {
	DescribeRestoreJob(ctx context.Context, restoreJobId string) (*pinecone.RestoreJob, error)
	ListRestoreJobs(ctx context.Context, in *pinecone.ListRestoreJobsParams) (*pinecone.RestoreJobList, error)
	CreateIndexFromBackup(ctx context.Context, in *pinecone.CreateIndexFromBackupParams) (*pinecone.CreateIndexFromBackupResponse, error)
}

func NewRestoreJobCmd() *cobra.Command {
	options := restoreJobCmdOptions{}
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Inspect restore jobs",
		Long:  restoreJobHelp,
		Example: help.Examples(`
			# Restore indexes from backups, creating restore jobs
			pc pinecone backup restore --id backup-123 --name restored-index
			pc pinecone backup restore --id backup-123 --name restored-index --tags env=prod,team=search --deletion-protection enabled

			# List restore jobs
			pc pinecone backup restore list

			# Describe a restore job
			pc pinecone backup restore describe --id rj-123
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			err := runRestoreJobCmd(ctx, pc, options)
			if err != nil {
				msg.FailMsg("Failed to create restore job: %s\n", err)
				exit.Error(err, "Failed to create restore job")
			}
		},
	}

	cmd.AddCommand(NewDescribeRestoreJobCmd())
	cmd.AddCommand(NewListRestoreJobsCmd())

	cmd.Flags().StringVarP(&options.backupId, "id", "i", "", "ID of the backup to restore from")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the index to create from the backup")
	cmd.Flags().StringVarP(&options.deletionProtection, "deletion-protection", "d", "", "whether to enable deletion protection on the new index (enabled|disabled)")
	cmd.Flags().StringToStringVarP(&options.tags, "tags", "t", map[string]string{}, "tags to apply to the new index")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runRestoreJobCmd(ctx context.Context, svc RestoreJobService, options restoreJobCmdOptions) error {
	if strings.TrimSpace(options.backupId) == "" {
		return pcio.Errorf("--id is required")
	}
	if strings.TrimSpace(options.name) == "" {
		return pcio.Errorf("--name is required")
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
		pcio.Println(text.IndentJSON(resp))
		return nil
	}

	msg.SuccessMsg("Restore job %s started for backup %s.\n", style.Emphasis(resp.RestoreJobId), style.Emphasis(options.backupId))
	msg.InfoMsg("Created index ID: %s\n", style.Emphasis(resp.IndexId))
	msg.InfoMsg("Use %s to monitor progress.\n", style.Code("pc pinecone backup restore describe --id "+resp.RestoreJobId))
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
		return nil, pcio.Errorf("invalid deletion-protection value %q, must be one of: enabled, disabled", input)
	}
}
