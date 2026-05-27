package backup

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	backupHelp = help.Long(`
		Manage backups for serverless indexes. A backup is a static copy of a serverless index
		that only consumes storage. It is a non-queryable representation of a set of records.
		You can create a backup of a serverless index, and you can create a new index from a backup.

		Use these commands to create, describe, list, and delete backups.

		See: https://docs.pinecone.io/guides/manage-data/backups-overview
	`)
)

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup",
		Short:   "Manage serverless index backups",
		Long:    backupHelp,
		GroupID: help.GROUP_INDEX_MANAGEMENT.ID,
		Example: help.Examples(`
			# Create a backup for a serverless index
			pc index backup create --index-name my-index --name daily-backup

			# List backups for a serverless index
			pc index backup list --index-name my-index

			# List all backups in the project
			pc index backup list

			# Describe a backup
			pc index backup describe --id backup-123

			# Delete a backup
			pc index backup delete --id backup-123
		`),
	}

	cmd.AddCommand(NewCreateBackupCmd())
	cmd.AddCommand(NewDescribeBackupCmd())
	cmd.AddCommand(NewListBackupsCmd())
	cmd.AddCommand(NewDeleteBackupCmd())

	return cmd
}
