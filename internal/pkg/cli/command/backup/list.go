package backup

import (
	"context"
	"fmt"
	"sort"

	backuppresenters "github.com/pinecone-io/cli/internal/pkg/utils/backup/presenters"
	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
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
				errorutil.HandleAPIError(err, cmd, args)
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
				backuppresenters.PrintBackupTable(backups.Data)
			}
		},
	}

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.indexName, "index", "i", "", "filter backups by index name")

	return cmd
}
