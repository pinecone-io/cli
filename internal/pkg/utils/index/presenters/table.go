package presenters

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/index"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// IndexDisplayData represents the unified display structure for index information
type IndexDisplayData struct {
	// Essential information
	Name          string
	Specification string
	VectorType    string
	Metric        string
	Dimension     string

	// State information (only for existing indexes)
	Status             string
	Host               string
	DeletionProtection string

	// Pod-specific information
	Environment string
	PodType     string
	Replicas    string
	ShardCount  string
	PodCount    string

	// Serverless-specific information
	CloudProvider string
	Region        string

	// Inference information
	Model              string
	EmbeddingDimension string
	FieldMap           string
	ReadParameters     string
	WriteParameters    string

	// Other information
	Tags string
}

// ConvertIndexToDisplayData converts a pinecone.Index to IndexDisplayData
func ConvertIndexToDisplayData(idx *pinecone.Index) *IndexDisplayData {
	data := &IndexDisplayData{}

	// Essential information
	data.Name = idx.Name
	data.VectorType = string(idx.VectorType)
	data.Metric = string(idx.Metric)
	if idx.Dimension != nil && *idx.Dimension > 0 {
		data.Dimension = fmt.Sprintf("%d", *idx.Dimension)
	}

	// Determine specification
	if idx.Spec.Serverless == nil {
		data.Specification = "pod"
	} else {
		data.Specification = "serverless"
	}

	// State information
	if idx.Status != nil {
		data.Status = string(idx.Status.State)
	}
	data.Host = idx.Host
	if idx.DeletionProtection == pinecone.DeletionProtectionEnabled {
		data.DeletionProtection = "enabled"
	} else {
		data.DeletionProtection = "disabled"
	}

	// Pod-specific information
	if idx.Spec.Pod != nil {
		data.Environment = idx.Spec.Pod.Environment
		data.PodType = idx.Spec.Pod.PodType
		data.Replicas = fmt.Sprintf("%d", idx.Spec.Pod.Replicas)
		data.ShardCount = fmt.Sprintf("%d", idx.Spec.Pod.ShardCount)
		data.PodCount = fmt.Sprintf("%d", idx.Spec.Pod.PodCount)
	}

	// Serverless-specific information
	if idx.Spec.Serverless != nil {
		data.CloudProvider = string(idx.Spec.Serverless.Cloud)
		data.Region = idx.Spec.Serverless.Region
	}

	// Inference information
	if idx.Embed != nil {
		data.Model = idx.Embed.Model
		if idx.Embed.Dimension != nil && *idx.Embed.Dimension > 0 {
			data.EmbeddingDimension = fmt.Sprintf("%d", *idx.Embed.Dimension)
		}

		// Format field map
		if idx.Embed.FieldMap != nil && len(*idx.Embed.FieldMap) > 0 {
			var fieldMapPairs []string
			for k, v := range *idx.Embed.FieldMap {
				fieldMapPairs = append(fieldMapPairs, fmt.Sprintf("%s=%v", k, v))
			}
			slices.Sort(fieldMapPairs)
			data.FieldMap = strings.Join(fieldMapPairs, ", ")
		}

		// Format read parameters
		if idx.Embed.ReadParameters != nil && len(*idx.Embed.ReadParameters) > 0 {
			var readParamsPairs []string
			for k, v := range *idx.Embed.ReadParameters {
				readParamsPairs = append(readParamsPairs, fmt.Sprintf("%s=%v", k, v))
			}
			slices.Sort(readParamsPairs)
			data.ReadParameters = strings.Join(readParamsPairs, ", ")
		}

		// Format write parameters
		if idx.Embed.WriteParameters != nil && len(*idx.Embed.WriteParameters) > 0 {
			var writeParamsPairs []string
			for k, v := range *idx.Embed.WriteParameters {
				writeParamsPairs = append(writeParamsPairs, fmt.Sprintf("%s=%v", k, v))
			}
			slices.Sort(writeParamsPairs)
			data.WriteParameters = strings.Join(writeParamsPairs, ", ")
		}
	}

	// Tags
	if idx.Tags != nil && len(*idx.Tags) > 0 {
		var tagStrings []string
		for key, value := range *idx.Tags {
			tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
		}
		slices.Sort(tagStrings)
		data.Tags = strings.Join(tagStrings, ", ")
	}

	return data
}

// ConvertCreateOptionsToDisplayData converts index.CreateOptions to IndexDisplayData
func ConvertCreateOptionsToDisplayData(config *index.CreateOptions) *IndexDisplayData {
	data := &IndexDisplayData{}

	// Essential information
	data.Name = formatValueWithInferred(config.Name.Value, config.Name.Inferred)
	data.VectorType = formatValueWithInferred(config.VectorType.Value, config.VectorType.Inferred)
	data.Metric = formatValueWithInferred(config.Metric.Value, config.Metric.Inferred)
	if config.Dimension.Value > 0 {
		data.Dimension = formatValueWithInferred(fmt.Sprintf("%d", config.Dimension.Value), config.Dimension.Inferred)
	}
	data.Model = formatValueWithInferred(config.Model.Value, config.Model.Inferred)

	// Determine specification
	spec, specInferred := config.GetSpecString()
	data.Specification = formatValueWithInferred(spec, specInferred)

	// Pod-specific information
	if config.GetSpec() == index.IndexSpecPod {
		data.Environment = formatValueWithInferred(config.Environment.Value, config.Environment.Inferred)
		data.PodType = formatValueWithInferred(config.PodType.Value, config.PodType.Inferred)
		data.Replicas = formatValueWithInferred(fmt.Sprintf("%d", config.Replicas.Value), config.Replicas.Inferred)
		data.ShardCount = formatValueWithInferred(fmt.Sprintf("%d", config.Shards.Value), config.Shards.Inferred)
		// Pod count not available in create options
	}

	// Serverless-specific information
	if config.GetSpec() == index.IndexSpecServerless {
		data.CloudProvider = formatValueWithInferred(config.Cloud.Value, config.Cloud.Inferred)
		data.Region = formatValueWithInferred(config.Region.Value, config.Region.Inferred)
	}

	// Format field map
	if len(config.FieldMap.Value) > 0 {
		var fieldMapPairs []string
		for k, v := range config.FieldMap.Value {
			fieldMapPairs = append(fieldMapPairs, fmt.Sprintf("%s=%v", k, v))
		}
		data.FieldMap = formatValueWithInferred(strings.Join(fieldMapPairs, ", "), config.FieldMap.Inferred)
	}

	// Format read parameters
	if len(config.ReadParameters.Value) > 0 {
		var readParamsPairs []string
		for k, v := range config.ReadParameters.Value {
			readParamsPairs = append(readParamsPairs, fmt.Sprintf("%s=%v", k, v))
		}
		data.ReadParameters = formatValueWithInferred(strings.Join(readParamsPairs, ", "), config.ReadParameters.Inferred)
	}

	// Format write parameters
	if len(config.WriteParameters.Value) > 0 {
		var writeParamsPairs []string
		for k, v := range config.WriteParameters.Value {
			writeParamsPairs = append(writeParamsPairs, fmt.Sprintf("%s=%v", k, v))
		}
		data.WriteParameters = formatValueWithInferred(strings.Join(writeParamsPairs, ", "), config.WriteParameters.Inferred)
	}

	// Deletion protection
	deletionProtection := config.DeletionProtection.Value
	if deletionProtection == "" {
		deletionProtection = "disabled"
	}
	data.DeletionProtection = formatValueWithInferred(deletionProtection, config.DeletionProtection.Inferred)

	// Tags
	if len(config.Tags.Value) > 0 {
		var tagStrings []string
		for key, value := range config.Tags.Value {
			tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
		}
		data.Tags = formatValueWithInferred(strings.Join(tagStrings, ", "), config.Tags.Inferred)
	}

	return data
}

// PrintIndexDisplayTable creates and renders a table for index display data
func PrintIndexDisplayTable(data *IndexDisplayData) {
	// Build rows for the table
	var rows []presenters.Row

	// Essential information
	rows = append(rows, presenters.Row{"Name", data.Name})
	rows = append(rows, presenters.Row{"Specification", data.Specification})
	rows = append(rows, presenters.Row{"Vector Type", data.VectorType})
	rows = append(rows, presenters.Row{"Metric", data.Metric})
	rows = append(rows, presenters.Row{"Dimension", data.Dimension})

	// Add spacing
	rows = append(rows, presenters.Row{"", ""})

	// State information (only show if we have status data)
	if data.Status != "" {
		rows = append(rows, presenters.Row{"Status", data.Status})
		rows = append(rows, presenters.Row{"Host URL", data.Host})
		rows = append(rows, presenters.Row{"Deletion Protection", data.DeletionProtection})
		rows = append(rows, presenters.Row{"", ""})
	}

	// Spec-specific information
	if strings.HasPrefix(data.Specification, "serverless") {
		rows = append(rows, presenters.Row{"Cloud Provider", data.CloudProvider})
		rows = append(rows, presenters.Row{"Region", data.Region})
	} else if strings.HasPrefix(data.Specification, "pod") {
		rows = append(rows, presenters.Row{"Environment", data.Environment})
		rows = append(rows, presenters.Row{"Pod Type", data.PodType})
		rows = append(rows, presenters.Row{"Replicas", data.Replicas})
		rows = append(rows, presenters.Row{"Shard Count", data.ShardCount})
		if data.PodCount != "" {
			rows = append(rows, presenters.Row{"Pod Count", data.PodCount})
		}
	}

	// Add spacing
	rows = append(rows, presenters.Row{"", ""})

	// Inference information (only show if we have model data)
	if data.Model != "" {
		rows = append(rows, presenters.Row{"Model", data.Model})
		if data.EmbeddingDimension != "" {
			rows = append(rows, presenters.Row{"Embedding Dimension", data.EmbeddingDimension})
		}
		if data.FieldMap != "" {
			rows = append(rows, presenters.Row{"Field Map", data.FieldMap})
		}
		if data.ReadParameters != "" {
			rows = append(rows, presenters.Row{"Read Parameters", data.ReadParameters})
		}
		if data.WriteParameters != "" {
			rows = append(rows, presenters.Row{"Write Parameters", data.WriteParameters})
		}
		rows = append(rows, presenters.Row{"", ""})
	}

	// Other information
	if data.DeletionProtection != "" && data.Status == "" {
		rows = append(rows, presenters.Row{"Deletion Protection", data.DeletionProtection})
	}

	if data.Tags != "" {
		rows = append(rows, presenters.Row{"Tags", data.Tags})
	}

	// Print each row with right-aligned first column and secondary text styling
	for _, row := range rows {
		if len(row) >= 2 {
			// Right align the first column content
			rightAlignedFirstCol := fmt.Sprintf("%20s", row[0])

			// Apply secondary text styling to the first column
			styledFirstCol := style.SecondaryTextStyle().Render(rightAlignedFirstCol)

			// Print the row
			rowText := fmt.Sprintf("%s  %s", styledFirstCol, row[1])
			fmt.Println(rowText)
		} else if len(row) == 1 && row[0] == "" {
			// Empty row for spacing
			fmt.Println()
		}
	}
	// Add spacing after the last row
	fmt.Println()
}

// PrintIndexTableWithIndexAttributesGroups creates and renders a table for index information with custom index attribute groups
func PrintIndexTableWithIndexAttributesGroups(indexes []*pinecone.Index, groups []IndexAttributesGroup) {
	// Filter out groups that have no meaningful data
	nonEmptyGroups := filterNonEmptyIndexAttributesGroups(indexes, groups)
	if len(nonEmptyGroups) == 0 {
		return
	}

	// Get columns for the non-empty groups
	columns := GetColumnsForIndexAttributesGroups(nonEmptyGroups)

	// Build table rows
	var rows []presenters.Row
	for _, idx := range indexes {
		values := ExtractValuesForIndexAttributesGroups(idx, nonEmptyGroups)
		rows = append(rows, presenters.Row(values))
	}

	// Use the table utility
	presenters.PrintTable(presenters.TableOptions{
		Columns: columns,
		Rows:    rows,
	})

	fmt.Println()

	// Add a note about full URLs if state info is shown
	hasStateGroup := false
	for _, group := range nonEmptyGroups {
		if group == IndexAttributesGroupState {
			hasStateGroup = true
			break
		}
	}
	if hasStateGroup && len(indexes) > 0 {
		hint := fmt.Sprintf("Use %s to see index details", style.Code("pc index describe <name>"))
		fmt.Println(style.Hint(hint))
	}
}

// PrintDescribeIndexTable creates and renders a table for index description with right-aligned first column and secondary text styling
func PrintDescribeIndexTable(idx *pinecone.Index) {
	// Print title
	fmt.Println(style.Heading("Index Configuration"))
	fmt.Println()

	// Convert to display data and print
	data := ConvertIndexToDisplayData(idx)
	PrintIndexDisplayTable(data)
}

// PrintIndexCreateConfigTable creates and renders a table for index creation configuration
func PrintIndexCreateConfigTable(config *index.CreateOptions) {
	fmt.Println(style.Heading("Index Configuration"))
	fmt.Println()

	// Convert to display data and print with inferred values
	data := ConvertCreateOptionsToDisplayData(config)
	PrintIndexDisplayTable(data)
}

// ColorizeState applies appropriate styling to index state
func ColorizeState(state pinecone.IndexStatusState) string {
	switch state {
	case pinecone.Ready:
		return style.SuccessStyle().Render(string(state))
	case pinecone.Initializing, pinecone.Terminating, pinecone.ScalingDown, pinecone.ScalingDownPodSize, pinecone.ScalingUp, pinecone.ScalingUpPodSize:
		return style.WarningStyle().Render(string(state))
	case pinecone.InitializationFailed:
		return style.ErrorStyle().Render(string(state))
	default:
		return string(state)
	}
}

// ColorizeDeletionProtection applies appropriate styling to deletion protection status
func ColorizeDeletionProtection(deletionProtection pinecone.DeletionProtection) string {
	if deletionProtection == pinecone.DeletionProtectionEnabled {
		return style.SuccessStyle().Render("enabled")
	}
	return style.ErrorStyle().Render("disabled")
}

// formatValueWithInferred formats a value with "(inferred)" indicator if the value was inferred
func formatValueWithInferred(value string, inferred bool) string {
	if inferred {
		return fmt.Sprintf("%s %s", value, style.SecondaryTextStyle().Render("(inferred)"))
	}
	return value
}
