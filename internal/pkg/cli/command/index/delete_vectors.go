package index

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/flags"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type deleteVectorsCmdOptions struct {
	indexName        string
	namespace        string
	ids              flags.StringList
	filter           flags.JSONObject
	deleteAllVectors bool
	json             bool
}

func NewDeleteVectorsCmd() *cobra.Command {
	options := deleteVectorsCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete-vectors",
		Short: "Delete vectors from an index",
		Example: help.Examples(`
			pc index delete-vectors --index-name my-index --namespace my-namespace --ids my-id
			pc index delete-vectors --index-name my-index --namespace my-namespace --all-vectors
			pc index delete-vectors --index-name my-index --namespace my-namespace --filter '{"genre": "classical"}'
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runDeleteVectorsCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to delete vectors from")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to delete vectors from")
	cmd.Flags().Var(&options.ids, "ids", "IDs of the vectors to delete")
	cmd.Flags().Var(&options.filter, "filter", "filter to delete the vectors with")
	cmd.Flags().BoolVar(&options.deleteAllVectors, "all-vectors", false, "delete all vectors from the namespace")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runDeleteVectorsCmd(ctx context.Context, options deleteVectorsCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	if options.ids == nil && options.filter == nil && !options.deleteAllVectors {
		msg.FailMsg("Either --ids, --filter, or --all-vectors must be provided")
		exit.ErrorMsg("Either --ids, --filter, or --all-vectors must be provided")
	}

	// Delete all vectors in namespace
	if options.deleteAllVectors {
		err = ic.DeleteAllVectorsInNamespace(ctx)
		if err != nil {
			msg.FailMsg("Failed to delete all vectors in namespace: %s", err)
			exit.Error(err, "Failed to delete all vectors in namespace")
		}
		if !options.json {
			msg.SuccessMsg("Deleted all vectors in namespace: %s", options.namespace)
		}
		return
	}

	// Delete vectors by ID
	if len(options.ids) > 0 {
		err = ic.DeleteVectorsById(ctx, options.ids)
		if err != nil {
			msg.FailMsg("Failed to delete vectors by IDs: %s", err)
			exit.Error(err, "Failed to delete vectors by IDs")
		}
		if !options.json {
			msg.SuccessMsg("Deleted vectors by IDs: %s", options.ids)
		}
		return
	}

	// Delete vectors by filter
	if options.filter != nil {
		filter, err := pinecone.NewMetadataFilter(options.filter)
		if err != nil {
			msg.FailMsg("Failed to create filter: %s", err)
			exit.Errorf(err, "Failed to create filter")
		}

		err = ic.DeleteVectorsByFilter(ctx, filter)
		if err != nil {
			msg.FailMsg("Failed to delete vectors by filter: %s", err)
			exit.Error(err, "Failed to delete vectors by filter")
		}
		if !options.json {
			msg.SuccessMsg("Deleted vectors by filter: %s", filter.String())
		}
		return
	}
}
