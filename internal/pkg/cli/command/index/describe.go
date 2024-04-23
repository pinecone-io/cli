package index

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/client"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/pinecone"
)

type DescribeCmdOptions struct {
	name string
	json bool
}

func NewDescribeCmd() *cobra.Command {
	options := DescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get configuration and status information for an index",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := client.NewPineconeClient()

			idx, err := pc.DescribeIndex(ctx, options.name)
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(idx)
			} else {
				printDescribeIndexTable(idx)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe")
	cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func printDescribeIndexTable(idx *pinecone.Index) {
	writer := tabwriter.NewWriter(os.Stdout, 12, 1, 4, ' ', 0)

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprintf(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", idx.Name)
	fmt.Fprintf(writer, "Dimension\t%d\n", idx.Dimension)
	fmt.Fprintf(writer, "Metric\t%s\n", string(idx.Metric))
	fmt.Fprintf(writer, "State\t%s\n", string(idx.Status.State))
	if idx.Status.Ready {
		fmt.Fprintf(writer, "Ready\t%s\n", style.StatusGreen("true"))
	} else {
		fmt.Fprintf(writer, "Ready\t%s\n", style.StatusRed("false"))
	}
	fmt.Fprintf(writer, "Host\t%s\n", style.Emphasis(idx.Host))

	var specType string
	if idx.Spec.Serverless == nil {
		specType = "pod"
		fmt.Fprintf(writer, "Spec\t%s\n", specType)
		fmt.Fprintf(writer, "Cloud\t%s\n", "")
		fmt.Fprintf(writer, "Region\t%s\n", "")
		fmt.Fprintf(writer, "Environment\t%s\n", idx.Spec.Pod.Environment)
	} else {
		specType = "serverless"
		fmt.Fprintf(writer, "Spec\t%s\n", specType)
		fmt.Fprintf(writer, "Cloud\t%s\n", idx.Spec.Serverless.Cloud)
		fmt.Fprintf(writer, "Region\t%s\n", idx.Spec.Serverless.Region)
		fmt.Fprintf(writer, "Environment\t%s\n", "")
	}

	writer.Flush()
}
