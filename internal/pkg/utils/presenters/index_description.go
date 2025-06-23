package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
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
	log.Debug().Str("name", idx.Name).Msg("Printing index description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", idx.Name)
	if idx.Dimension != nil {
		pcio.Fprintf(writer, "Dimension\t%d\n", *idx.Dimension)
	} else {
		pcio.Fprintf(writer, "Dimension\tnil\n")
	}
	pcio.Fprintf(writer, "Metric\t%s\n", string(idx.Metric))
	pcio.Fprintf(writer, "Deletion Protection\t%s\n", ColorizeDeletionProtection(idx.DeletionProtection))
	pcio.Fprintf(writer, "\t\n")
	pcio.Fprintf(writer, "State\t%s\n", ColorizeState(idx.Status.State))
	pcio.Fprintf(writer, "Ready\t%s\n", ColorizeBool(idx.Status.Ready))
	pcio.Fprintf(writer, "Host\t%s\n", style.Emphasis(idx.Host))
	pcio.Fprintf(writer, "\t\n")

	var specType string
	if idx.Spec.Serverless == nil {
		specType = "pod"
		pcio.Fprintf(writer, "Spec\t%s\n", specType)
		pcio.Fprintf(writer, "Environment\t%s\n", idx.Spec.Pod.Environment)
		pcio.Fprintf(writer, "PodType\t%s\n", idx.Spec.Pod.PodType)
		pcio.Fprintf(writer, "Replicas\t%d\n", idx.Spec.Pod.Replicas)
		pcio.Fprintf(writer, "ShardCount\t%d\n", idx.Spec.Pod.ShardCount)
		pcio.Fprintf(writer, "PodCount\t%d\n", idx.Spec.Pod.PodCount)
		pcio.Fprintf(writer, "MetadataConfig\t%s\n", text.InlineJSON(idx.Spec.Pod.MetadataConfig))
	} else {
		specType = "serverless"
		pcio.Fprintf(writer, "Spec\t%s\n", specType)
		pcio.Fprintf(writer, "Cloud\t%s\n", idx.Spec.Serverless.Cloud)
		pcio.Fprintf(writer, "Region\t%s\n", idx.Spec.Serverless.Region)
	}

	writer.Flush()
}
