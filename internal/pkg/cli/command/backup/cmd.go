package backup

import (
	"github.com/pinecone-io/cli/internal/pkg/cli/command/backup/restore"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	backupHelp = help.Long(`
		Manage backups for serverless indexes. A backup is a static copy of a serverless index
		that only consumes storage. It is a non-queryable representation of a set of records.
		You can create a backup of a serverless index, and you can create a new index from a backup.

		Use these commands to create, describe, list, and delete backups, or to
		create a new index from an existing backup. You can also describe and list restore jobs.

		See: https://docs.pinecone.io/guides/manage-data/backups-overview
	`)
)

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Manage serverless index backups",
		Long:  backupHelp,
		Example: help.Examples(`
			pc pinecone backup create --index-name my-index --name daily-backup
			pc pinecone backup list --index-name my-index
			pc pinecone backup create-index --id backup-123 --name restored-index
		`),
	}

	cmd.AddCommand(NewCreateBackupCmd())
	cmd.AddCommand(NewDescribeBackupCmd())
	cmd.AddCommand(NewListBackupsCmd())
	cmd.AddCommand(NewDeleteBackupCmd())

	cmd.AddCommand(restore.NewRestoreJobCmd())

	return cmd
}
