package index

import (
	"context"
	"fmt"

	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/index"
	indexpresenters "github.com/pinecone-io/cli/internal/pkg/utils/index/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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
	json               bool
}

func NewConfigureIndexCmd() *cobra.Command {
	options := configureIndexOptions{}

	cmd := &cobra.Command{
		Use:          "configure <name>",
		Short:        "Configure an existing index with the specified configuration",
		Example:      "",
		Args:         index.ValidateIndexNameArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			options.name = args[0]
			runConfigureIndexCmd(options, cmd, args)
		},
	}

	// Optional flags
	cmd.Flags().StringVarP(&options.podType, "pod_type", "t", "", "type of pod to use, can only upgrade when configuring")
	cmd.Flags().Int32VarP(&options.replicas, "replicas", "r", 0, "replicas of the index to configure")
	cmd.Flags().StringVarP(&options.deletionProtection, "deletion_protection", "p", "", "enable or disable deletion protection for the index")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func runConfigureIndexCmd(options configureIndexOptions, cmd *cobra.Command, args []string) {
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	idx, err := pc.ConfigureIndex(ctx, options.name, pinecone.ConfigureIndexParams{
		PodType:            options.podType,
		Replicas:           options.replicas,
		DeletionProtection: pinecone.DeletionProtection(options.deletionProtection),
	})
	if err != nil {
		errorutil.HandleAPIError(err, cmd, args)
		exit.Error(err)
	}

	if options.json {
		json := text.IndentJSON(idx)
		pcio.Println(json)
		return
	}

	msg.SuccessMsg("Index %s configured successfully.", style.ResourceName(idx.Name))

	indexpresenters.PrintDescribeIndexTable(idx)

	describeCommand := pcio.Sprintf("pc index describe %s", idx.Name)
	hint := fmt.Sprintf("Run %s at any time to check the status. \n\n", style.Code(describeCommand))
	pcio.Println(style.Hint(hint))
}
