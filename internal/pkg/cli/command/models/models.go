package models

import (
	"context"
	_ "embed"
	"fmt"
	"sort"

	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/models/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type ListModelsCmdOptions struct {
	json bool
}

func NewModelsCmd() *cobra.Command {
	options := ListModelsCmdOptions{}

	cmd := &cobra.Command{
		Use:   "models",
		Short: "List the models hosted on Pinecone",
		Run: func(cmd *cobra.Command, args []string) {
			pc := sdk.NewPineconeClient()
			ctx := context.Background()

			embed := "embed"
			embedModels, err := pc.Inference.ListModels(ctx, &pinecone.ListModelsParams{Type: &embed})
			if err != nil {
				errorutil.HandleIndexAPIError(err, cmd, []string{})
				exit.Error(err)
			}

			if embedModels == nil || embedModels.Models == nil || len(*embedModels.Models) == 0 {
				fmt.Println("No models found.")
				return
			}

			models := *embedModels.Models

			// Sort results alphabetically by model name
			sort.SliceStable(models, func(i, j int) bool {
				return models[i].Model < models[j].Model
			})

			if options.json {
				// Use fmt for data output - should not be suppressed by -q flag
				json := text.IndentJSON(models)
				fmt.Println(json)
			} else {
				// Show models in table format
				presenters.PrintModelsTable(models)
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
