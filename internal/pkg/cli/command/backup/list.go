package backup

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

type ListBackupsCmdOptions struct {
	json      bool
	indexName string
}

func NewListBackupsCmd() *cobra.Command {
	options := ListBackupsCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "See the list of backups in your project",
		Run: func(cmd *cobra.Command, args []string) {
			pc := sdk.NewPineconeClient()
			ctx := context.Background()

			params := &pinecone.ListBackupsParams{}
			if options.indexName != "" {
				params.IndexName = &options.indexName
			}

			backups, err := pc.ListBackups(ctx, params)
			if err != nil {
				msg.FailMsg("Failed to list backups: %s\n", err)
				exit.Error(err)
			}

			// Sort results alphabetically by name
			sort.SliceStable(backups.Data, func(i, j int) bool {
				nameI := "unnamed"
				if backups.Data[i].Name != nil {
					nameI = *backups.Data[i].Name
				}
				nameJ := "unnamed"
				if backups.Data[j].Name != nil {
					nameJ = *backups.Data[j].Name
				}
				return nameI < nameJ
			})

			if options.json {
				json := text.IndentJSON(backups)
				fmt.Println(json)
			} else {
				printTable(backups.Data)
			}
		},
	}

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.indexName, "index", "i", "", "filter backups by index name")

	return cmd
}

func printTable(backups []*pinecone.Backup) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "ID", "INDEX", "STATUS", "CREATED", "SIZE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	for _, backup := range backups {
		created := "-"
		if backup.CreatedAt != nil {
			created = *backup.CreatedAt
		}

		size := "-"
		if backup.SizeBytes != nil {
			size = fmt.Sprintf("%d", *backup.SizeBytes)
		}

		backupName := "unnamed"
		if backup.Name != nil {
			backupName = *backup.Name
		}
		values := []string{
			backupName,
			backup.BackupId,
			backup.SourceIndexName,
			backup.Status,
			created,
			size,
		}
		fmt.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
