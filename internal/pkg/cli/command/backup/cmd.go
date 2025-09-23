package backup

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

var backupHelpText = text.WordWrap(`
A backup is a snapshot of a serverless index that can be used to restore 
the index to a previous state or create a new serverless index. Backups are 
useful for disaster recovery, data migration, and creating development 
environments from production data.
`, 80)

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup <command>",
		Short: "Work with backups",
		Long:  backupHelpText,
		Example: heredoc.Doc(`
			$ pc backup list
			$ pc backup create --name my-backup --source my-serverless-index
			$ pc backup describe backup-123
			$ pc backup restore --index restored-index --backup-id backup-123
			$ pc backup delete backup-123
		`),
		GroupID: help.GROUP_VECTORDB.ID,
	}

	cmd.AddCommand(NewCreateBackupCmd())
	cmd.AddCommand(NewListBackupsCmd())
	cmd.AddCommand(NewDescribeBackupCmd())
	cmd.AddCommand(NewDeleteBackupCmd())
	cmd.AddCommand(NewRestoreBackupCmd())

	return cmd
}
