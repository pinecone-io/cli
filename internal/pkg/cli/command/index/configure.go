package index

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type configureIndexOptions struct {
	name               string
	podType            string
	replicas           int32
	deletionProtection string

	json bool
}

func NewConfigureIndexCmd() *cobra.Command {
	options := configureIndexOptions{}

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure an existing index with the specified configuration",
		Example: help.Examples(`
			pc index configure --name my-index --deletion-protection enabled
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runConfigureIndexCmd(options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to configure")

	// Optional flags
	cmd.Flags().StringVarP(&options.podType, "pod_type", "t", "", "type of pod to use, can only upgrade when configuring")
	cmd.Flags().Int32VarP(&options.replicas, "replicas", "r", 0, "replicas of the index to configure")
	cmd.Flags().StringVarP(&options.deletionProtection, "deletion_protection", "p", "", "enable or disable deletion protection for the index")

	return cmd
}

func runConfigureIndexCmd(options configureIndexOptions) {
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	idx, err := pc.ConfigureIndex(ctx, options.name, pinecone.ConfigureIndexParams{
		PodType:            options.podType,
		Replicas:           options.replicas,
		DeletionProtection: pinecone.DeletionProtection(options.deletionProtection),
	})
	if err != nil {
		msg.FailMsg("Failed to configure index %s: %+v\n", style.Emphasis(options.name), err)
		exit.Error(err)
	}
	if options.json {
		json := text.IndentJSON(idx)
		pcio.Println(json)
		return
	}

	describeCommand := pcio.Sprintf("pc index describe --name %s", idx.Name)
	msg.SuccessMsg("Index %s configured successfully. Run %s to check status. \n\n", style.Emphasis(idx.Name), style.Code(describeCommand))
	presenters.PrintDescribeIndexTable(idx)
}
