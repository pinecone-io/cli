package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

func ColorizeState(state pinecone.IndexStatusState) string {
	switch state {
	case pinecone.Ready:
		return style.StatusGreen(string(state))
	case pinecone.Initializing, pinecone.Terminating, pinecone.ScalingDown, pinecone.ScalingDownPodSize, pinecone.ScalingUp, pinecone.ScalingUpPodSize:
		return style.StatusYellow(string(state))
	case pinecone.InitializationFailed:
		return style.StatusRed(string(state))
	default:
		return string(state)
	}
}

func PrintDescribeIndexTable(idx *pinecone.Index) {
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", idx.Name)
	fmt.Fprintf(writer, "Dimension\t%d\n", idx.Dimension)
	fmt.Fprintf(writer, "Metric\t%s\n", string(idx.Metric))
	fmt.Fprintf(writer, "\t\n")
	fmt.Fprintf(writer, "State\t%s\n", ColorizeState(idx.Status.State))
	fmt.Fprintf(writer, "Ready\t%s\n", ColorizeBool(idx.Status.Ready))
	fmt.Fprintf(writer, "Host\t%s\n", style.Emphasis(idx.Host))
	fmt.Fprintf(writer, "\t\n")

	var specType string
	if idx.Spec.Serverless == nil {
		specType = "pod"
		fmt.Fprintf(writer, "Spec\t%s\n", specType)
		fmt.Fprintf(writer, "Environment\t%s\n", idx.Spec.Pod.Environment)
		fmt.Fprintf(writer, "PodType\t%s\n", idx.Spec.Pod.PodType)
		fmt.Fprintf(writer, "Replicas\t%d\n", idx.Spec.Pod.Replicas)
		fmt.Fprintf(writer, "ShardCount\t%d\n", idx.Spec.Pod.ShardCount)
		fmt.Fprintf(writer, "PodCount\t%d\n", idx.Spec.Pod.PodCount)
		fmt.Fprintf(writer, "MetadataConfig\t%s\n", text.InlineJSON(idx.Spec.Pod.MetadataConfig))
	} else {
		specType = "serverless"
		fmt.Fprintf(writer, "Spec\t%s\n", specType)
		fmt.Fprintf(writer, "Cloud\t%s\n", idx.Spec.Serverless.Cloud)
		fmt.Fprintf(writer, "Region\t%s\n", idx.Spec.Serverless.Region)
	}

	writer.Flush()
}
