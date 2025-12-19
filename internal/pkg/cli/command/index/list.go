package index

import (
	"sort"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type listIndexCmdOptions struct {
	json bool
	wide bool
}

func NewListCmd() *cobra.Command {
	options := listIndexCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all indexes in the target project",
		Example: help.Examples(`
			pc index list
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			idxs, err := pc.ListIndexes(ctx)
			if err != nil {
				msg.FailMsg("Failed to list indexes: %s\n", err)
				exit.Error(err, "Failed to list indexes")
			}

			// Sort results alphabetically by name
			sort.SliceStable(idxs, func(i, j int) bool {
				return idxs[i].Name < idxs[j].Name
			})

			if options.json {
				json := text.IndentJSON(idxs)
				pcio.Println(json)
			} else {
				printTable(idxs, options.wide)
			}
		},
	}

	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON, includes full index details")
	cmd.Flags().BoolVarP(&options.wide, "wide", "w", false, "Show additional columns (host, embed, tags)")

	return cmd
}

// printTable prints the index list in a table format
func printTable(idxs []*pinecone.Index, wide bool) {
	writer := presenters.NewTabWriter()

	columns := []string{"NAME", "STATUS", "SPEC", "CLOUD/REGION", "METRIC", "DIMENSION", "READ CAPACITY", "HOST"}
	if wide {
		columns = append(columns, "EMBED", "TAGS")
	}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, idx := range idxs {
		status := "-"
		if idx.Status != nil {
			status = string(idx.Status.State)
		}

		dimension := "nil"
		if idx.Dimension != nil {
			dimension = pcio.Sprintf("%d", *idx.Dimension)
		}

		spec := formatSpec(idx.Spec)
		cloudRegion := formatCloudRegion(idx.Spec)
		readCapacity := formatReadCapacity(idx.Spec)
		embed := formatEmbed(idx.Embed)
		tags := formatTags(idx.Tags)
		host := formatHost(idx.Host, wide)

		values := []string{idx.Name, status, spec, cloudRegion, string(idx.Metric), dimension, readCapacity, host}
		if wide {
			values = append(values, embed, tags)
		}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}

	if !wide {
		pcio.Fprint(writer, "\nUse --wide to show host/embed/tags, or --json for full details.\n")
	}
	writer.Flush()
}

// formatSpec formats the index spec as "serverless", "byoc", or "pod"
func formatSpec(spec *pinecone.IndexSpec) string {
	switch {
	case spec == nil:
		return "-"
	case spec.Serverless != nil:
		return "serverless"
	case spec.BYOC != nil:
		return "byoc"
	case spec.Pod != nil:
		return "pod"
	default:
		return "-"
	}
}

// formatCloudRegion formats the cloud and region as "cloud/region"
func formatCloudRegion(spec *pinecone.IndexSpec) string {
	if spec == nil || spec.Serverless == nil {
		return "-"
	}
	cloud := string(spec.Serverless.Cloud)
	region := spec.Serverless.Region
	if cloud == "" && region == "" {
		return "-"
	}
	if cloud == "" {
		return region
	}
	if region == "" {
		return cloud
	}
	return cloud + "/" + region
}

// formatReadCapacity formats the read capacity as "OnDemand" or "Dedicated"
func formatReadCapacity(spec *pinecone.IndexSpec) string {
	if spec == nil || spec.Serverless == nil || spec.Serverless.ReadCapacity == nil {
		return "-"
	}
	rc := spec.Serverless.ReadCapacity
	switch {
	case rc.OnDemand != nil:
		return "OnDemand"
	case rc.Dedicated != nil:
		return "Dedicated"
	default:
		return "-"
	}
}

// formatEmbed formats the embed model name
func formatEmbed(embed *pinecone.IndexEmbed) string {
	if embed == nil {
		return "-"
	}
	if embed.Model == "" {
		return "-"
	}
	return embed.Model
}

// formatTags formats the tags as a comma-separated list
func formatTags(tags *pinecone.IndexTags) string {
	if tags == nil {
		return "-"
	}

	keys := make([]string, 0, len(*tags))
	for k := range *tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if len(keys) == 0 {
		return "-"
	}

	return strings.Join(keys, ",")
}

// formatHost shortens the host when not in --wide mode
func formatHost(host string, wide bool) string {
	if host == "" {
		return "-"
	}
	if wide {
		return host
	}
	const maxLen = 30
	return midTruncate(host, maxLen)
}

// midTruncate shortens s to maxLen by keeping the start and end and inserting "...".
func midTruncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	// Reserve space for ellipsis.
	keep := maxLen - 3
	if keep <= 0 {
		return "..."
	}
	front := keep / 2
	back := keep - front

	return string(runes[:front]) + "..." + string(runes[len(runes)-back:])
}
