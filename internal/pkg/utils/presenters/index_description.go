package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
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

func ColorizeDeletionProtection(deletionProtection pinecone.DeletionProtection) string {
	if deletionProtection == pinecone.DeletionProtectionEnabled {
		return style.StatusGreen("enabled")
	}
	return style.StatusRed("disabled")
}

func PrintDescribeIndexTable(idx *pinecone.Index) {
	writer := NewTabWriter()
	if idx == nil {
		PrintEmptyState(writer, "index details")
		return
	}

	log.Debug().Str("name", idx.Name).Msg("Printing index description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", idx.Name)
	fmt.Fprintf(writer, "Dimension\t%s\n", DisplayOrNone(idx.Dimension))
	fmt.Fprintf(writer, "Metric\t%s\n", string(idx.Metric))
	fmt.Fprintf(writer, "Deletion Protection\t%s\n", ColorizeDeletionProtection(idx.DeletionProtection))
	fmt.Fprintf(writer, "Vector Type\t%s\n", DisplayOrNone(idx.VectorType))
	fmt.Fprintf(writer, "\n")

	stateVal := "<none>"
	readyVal := "<none>"
	if idx.Status != nil {
		stateVal = ColorizeState(idx.Status.State)
		readyVal = ColorizeBool(idx.Status.Ready)
	}
	fmt.Fprintf(writer, "State\t%s\n", stateVal)
	fmt.Fprintf(writer, "Ready\t%s\n", readyVal)
	fmt.Fprintf(writer, "Host\t%s\n", style.Emphasis(idx.Host))
	fmt.Fprintf(writer, "Private Host\t%s\n", DisplayOrNone(idx.PrivateHost))
	fmt.Fprintf(writer, "\n")

	switch {
	case idx.Spec == nil: // nil spec
		fmt.Fprintf(writer, "Spec\t%s\n", "<none>")
	case idx.Spec.Serverless != nil: // serverless spec
		fmt.Fprintf(writer, "Spec\t%s\n", "serverless")
		serverless := idx.Spec.Serverless
		if serverless != nil {
			fmt.Fprintf(writer, "Cloud\t%s\n", serverless.Cloud)
			fmt.Fprintf(writer, "Region\t%s\n", serverless.Region)
			fmt.Fprintf(writer, "Source Collection\t%s\n", DisplayOrNone(serverless.SourceCollection))
			schemaVal := "<none>"
			if serverless.Schema != nil {
				schemaVal = text.InlineJSON(serverless.Schema)
			}
			fmt.Fprintf(writer, "Schema\t%s\n", schemaVal)
			readCapacityVal := "<none>"
			if serverless.ReadCapacity != nil {
				readCapacityVal = text.InlineJSON(serverless.ReadCapacity)
			}
			fmt.Fprintf(writer, "Read Capacity\t%s\n", readCapacityVal)
		} else {
			fmt.Fprintf(writer, "Cloud\t%s\n", "<none>")
			fmt.Fprintf(writer, "Region\t%s\n", "<none>")
			fmt.Fprintf(writer, "Source Collection\t%s\n", "<none>")
		}
	case idx.Spec.BYOC != nil: // BYOC spec
		fmt.Fprintf(writer, "Spec\t%s\n", "byoc")
		byoc := idx.Spec.BYOC
		fmt.Fprintf(writer, "Environment\t%s\n", byoc.Environment)
		schemaVal := "<none>"
		if byoc.Schema != nil {
			schemaVal = text.InlineJSON(byoc.Schema)
		}
		fmt.Fprintf(writer, "Schema\t%s\n", schemaVal)
	case idx.Spec.Pod != nil: // pod spec
		fmt.Fprintf(writer, "Spec\t%s\n", "pod")
		pod := idx.Spec.Pod
		if pod != nil {
			fmt.Fprintf(writer, "Environment\t%s\n", pod.Environment)
			fmt.Fprintf(writer, "PodType\t%s\n", pod.PodType)
			fmt.Fprintf(writer, "Replicas\t%d\n", pod.Replicas)
			fmt.Fprintf(writer, "ShardCount\t%d\n", pod.ShardCount)
			fmt.Fprintf(writer, "PodCount\t%d\n", pod.PodCount)
			metadataConfig := "<none>"
			if pod.MetadataConfig != nil {
				metadataConfig = text.InlineJSON(pod.MetadataConfig)
			}
			fmt.Fprintf(writer, "MetadataConfig\t%s\n", metadataConfig)
			fmt.Fprintf(writer, "Source Collection\t%s\n", DisplayOrNone(pod.SourceCollection))
		} else {
			fmt.Fprintf(writer, "Environment\t%s\n", "<none>")
			fmt.Fprintf(writer, "PodType\t%s\n", "<none>")
			fmt.Fprintf(writer, "Replicas\t%s\n", "<none>")
			fmt.Fprintf(writer, "ShardCount\t%s\n", "<none>")
			fmt.Fprintf(writer, "PodCount\t%s\n", "<none>")
			fmt.Fprintf(writer, "MetadataConfig\t%s\n", "<none>")
			fmt.Fprintf(writer, "Source Collection\t%s\n", "<none>")
		}
	default: // unknown spec
		fmt.Fprintf(writer, "Spec\t%s\n", "unknown")
	}
	fmt.Fprintf(writer, "\n")

	if idx.Embed != nil {
		fmt.Fprintf(writer, "Model\t%s\n", idx.Embed.Model)
		fmt.Fprintf(writer, "Field Map\t%s\n", text.InlineJSON(idx.Embed.FieldMap))
		fmt.Fprintf(writer, "Read Parameters\t%s\n", text.InlineJSON(idx.Embed.ReadParameters))
		fmt.Fprintf(writer, "Write Parameters\t%s\n", text.InlineJSON(idx.Embed.WriteParameters))
		fmt.Fprintf(writer, "\n")
	}

	if idx.Tags != nil {
		fmt.Fprintf(writer, "Tags\t%s\n", text.InlineJSON(idx.Tags))
	}

	writer.Flush()
}
