package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", idx.Name)
	pcio.Fprintf(writer, "Dimension\t%s\n", DisplayOrNone(idx.Dimension))
	pcio.Fprintf(writer, "Metric\t%s\n", string(idx.Metric))
	pcio.Fprintf(writer, "Deletion Protection\t%s\n", ColorizeDeletionProtection(idx.DeletionProtection))
	pcio.Fprintf(writer, "Vector Type\t%s\n", DisplayOrNone(idx.VectorType))
	pcio.Fprintf(writer, "\t\n")
	stateVal := "<none>"
	readyVal := "<none>"
	if idx.Status != nil {
		stateVal = ColorizeState(idx.Status.State)
		readyVal = ColorizeBool(idx.Status.Ready)
	}
	pcio.Fprintf(writer, "State\t%s\n", stateVal)
	pcio.Fprintf(writer, "Ready\t%s\n", readyVal)
	pcio.Fprintf(writer, "Host\t%s\n", style.Emphasis(idx.Host))
	pcio.Fprintf(writer, "Private Host\t%s\n", DisplayOrNone(idx.PrivateHost))
	pcio.Fprintf(writer, "\t\n")

	switch {
	case idx.Spec == nil: // nil spec
		pcio.Fprintf(writer, "Spec\t%s\n", "<none>")
	case idx.Spec.Serverless != nil: // serverless spec
		pcio.Fprintf(writer, "Spec\t%s\n", "serverless")
		serverless := idx.Spec.Serverless
		if serverless != nil {
			pcio.Fprintf(writer, "Cloud\t%s\n", serverless.Cloud)
			pcio.Fprintf(writer, "Region\t%s\n", serverless.Region)
			pcio.Fprintf(writer, "Source Collection\t%s\n", DisplayOrNone(serverless.SourceCollection))
			schemaVal := "<none>"
			if serverless.Schema != nil {
				schemaVal = text.InlineJSON(serverless.Schema)
			}
			pcio.Fprintf(writer, "Schema\t%s\n", schemaVal)
			readCapacityVal := "<none>"
			if serverless.ReadCapacity != nil {
				readCapacityVal = text.InlineJSON(serverless.ReadCapacity)
			}
			pcio.Fprintf(writer, "Read Capacity\t%s\n", readCapacityVal)
		} else {
			pcio.Fprintf(writer, "Cloud\t%s\n", "<none>")
			pcio.Fprintf(writer, "Region\t%s\n", "<none>")
			pcio.Fprintf(writer, "Source Collection\t%s\n", "<none>")
		}
	case idx.Spec.BYOC != nil: // BYOC spec
		pcio.Fprintf(writer, "Spec\t%s\n", "byoc")
		byoc := idx.Spec.BYOC
		pcio.Fprintf(writer, "Environment\t%s\n", byoc.Environment)
		schemaVal := "<none>"
		if byoc.Schema != nil {
			schemaVal = text.InlineJSON(byoc.Schema)
		}
		pcio.Fprintf(writer, "Schema\t%s\n", schemaVal)
	case idx.Spec.Pod != nil: // pod spec
		pcio.Fprintf(writer, "Spec\t%s\n", "pod")
		pod := idx.Spec.Pod
		if pod != nil {
			pcio.Fprintf(writer, "Environment\t%s\n", pod.Environment)
			pcio.Fprintf(writer, "PodType\t%s\n", pod.PodType)
			pcio.Fprintf(writer, "Replicas\t%d\n", pod.Replicas)
			pcio.Fprintf(writer, "ShardCount\t%d\n", pod.ShardCount)
			pcio.Fprintf(writer, "PodCount\t%d\n", pod.PodCount)
			metadataConfig := "<none>"
			if pod.MetadataConfig != nil {
				metadataConfig = text.InlineJSON(pod.MetadataConfig)
			}
			pcio.Fprintf(writer, "MetadataConfig\t%s\n", metadataConfig)
			pcio.Fprintf(writer, "Source Collection\t%s\n", DisplayOrNone(pod.SourceCollection))
		} else {
			pcio.Fprintf(writer, "Environment\t%s\n", "<none>")
			pcio.Fprintf(writer, "PodType\t%s\n", "<none>")
			pcio.Fprintf(writer, "Replicas\t%s\n", "<none>")
			pcio.Fprintf(writer, "ShardCount\t%s\n", "<none>")
			pcio.Fprintf(writer, "PodCount\t%s\n", "<none>")
			pcio.Fprintf(writer, "MetadataConfig\t%s\n", "<none>")
			pcio.Fprintf(writer, "Source Collection\t%s\n", "<none>")
		}
	default: // unknown spec
		pcio.Fprintf(writer, "Spec\t%s\n", "unknown")
	}
	pcio.Fprintf(writer, "\t\n")

	if idx.Embed != nil {
		pcio.Fprintf(writer, "Model\t%s\n", idx.Embed.Model)
		pcio.Fprintf(writer, "Field Map\t%s\n", text.InlineJSON(idx.Embed.FieldMap))
		pcio.Fprintf(writer, "Read Parameters\t%s\n", text.InlineJSON(idx.Embed.ReadParameters))
		pcio.Fprintf(writer, "Write Parameters\t%s\n", text.InlineJSON(idx.Embed.WriteParameters))
	}

	if idx.Tags != nil {
		pcio.Fprintf(writer, "Tags\t%s\n", text.InlineJSON(idx.Tags))
	}

	writer.Flush()
}
