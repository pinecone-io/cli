package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// IndexAttributesGroup represents the available attribute groups for index display
type IndexAttributesGroup string

const (
	IndexAttributesGroupEssential      IndexAttributesGroup = "essential"
	IndexAttributesGroupState          IndexAttributesGroup = "state"
	IndexAttributesGroupPodSpec        IndexAttributesGroup = "pod_spec"
	IndexAttributesGroupServerlessSpec IndexAttributesGroup = "serverless_spec"
	IndexAttributesGroupInference      IndexAttributesGroup = "inference"
	IndexAttributesGroupOther          IndexAttributesGroup = "other"
)

// AllIndexAttributesGroups returns all available index attribute groups
func AllIndexAttributesGroups() []IndexAttributesGroup {
	return []IndexAttributesGroup{
		IndexAttributesGroupEssential,
		IndexAttributesGroupState,
		IndexAttributesGroupPodSpec,
		IndexAttributesGroupServerlessSpec,
		IndexAttributesGroupInference,
		IndexAttributesGroupOther,
	}
}

// IndexColumn represents a table column with both short and full names
type IndexColumn struct {
	ShortTitle string
	FullTitle  string
	Width      int
}

// ColumnGroup represents a group of columns with both short and full names
type ColumnGroup struct {
	Name    string
	Columns []IndexColumn
}

// IndexColumnGroups defines the available column groups for index tables
// Each group represents a logical set of related index properties that can be displayed together
var IndexColumnGroups = struct {
	Essential      ColumnGroup // Basic index information (name, spec, type, metric, dimension)
	State          ColumnGroup // Runtime state information (status, host, protection)
	PodSpec        ColumnGroup // Pod-specific configuration (environment, pod type, replicas, etc.)
	ServerlessSpec ColumnGroup // Serverless-specific configuration (cloud, region)
	Inference      ColumnGroup // Inference/embedding model information
	Other          ColumnGroup // Other information (tags, custom fields, etc.)
}{
	Essential: ColumnGroup{
		Name: "essential",
		Columns: []IndexColumn{
			{ShortTitle: "NAME", FullTitle: "Name", Width: 20},
			{ShortTitle: "SPEC", FullTitle: "Specification", Width: 12},
			{ShortTitle: "TYPE", FullTitle: "Vector Type", Width: 8},
			{ShortTitle: "METRIC", FullTitle: "Metric", Width: 8},
			{ShortTitle: "DIM", FullTitle: "Dimension", Width: 8},
		},
	},
	State: ColumnGroup{
		Name: "state",
		Columns: []IndexColumn{
			{ShortTitle: "STATUS", FullTitle: "Status", Width: 10},
			{ShortTitle: "HOST", FullTitle: "Host URL", Width: 60},
			{ShortTitle: "PROT", FullTitle: "Deletion Protection", Width: 8},
		},
	},
	PodSpec: ColumnGroup{
		Name: "pod_spec",
		Columns: []IndexColumn{
			{ShortTitle: "ENV", FullTitle: "Environment", Width: 12},
			{ShortTitle: "POD_TYPE", FullTitle: "Pod Type", Width: 12},
			{ShortTitle: "REPLICAS", FullTitle: "Replicas", Width: 8},
			{ShortTitle: "SHARDS", FullTitle: "Shard Count", Width: 8},
			{ShortTitle: "PODS", FullTitle: "Pod Count", Width: 8},
		},
	},
	ServerlessSpec: ColumnGroup{
		Name: "serverless_spec",
		Columns: []IndexColumn{
			{ShortTitle: "CLOUD", FullTitle: "Cloud Provider", Width: 12},
			{ShortTitle: "REGION", FullTitle: "Region", Width: 15},
		},
	},
	Inference: ColumnGroup{
		Name: "inference",
		Columns: []IndexColumn{
			{ShortTitle: "MODEL", FullTitle: "Model", Width: 25},
			{ShortTitle: "EMBED DIM", FullTitle: "Embedding Dimension", Width: 10},
			{ShortTitle: "FIELD MAP", FullTitle: "Field Map", Width: 20},
			{ShortTitle: "READ PARAMS", FullTitle: "Read Parameters", Width: 20},
			{ShortTitle: "WRITE PARAMS", FullTitle: "Write Parameters", Width: 20},
		},
	},
	Other: ColumnGroup{
		Name: "other",
		Columns: []IndexColumn{
			{ShortTitle: "TAGS", FullTitle: "Tags", Width: 30},
		},
	},
}

// GetColumnsForIndexAttributesGroups returns columns for the specified index attribute groups (using short names for horizontal tables)
func GetColumnsForIndexAttributesGroups(groups []IndexAttributesGroup) []Column {
	var columns []Column
	for _, group := range groups {
		switch group {
		case IndexAttributesGroupEssential:
			for _, col := range IndexColumnGroups.Essential.Columns {
				columns = append(columns, Column{Title: col.ShortTitle, Width: col.Width})
			}
		case IndexAttributesGroupState:
			for _, col := range IndexColumnGroups.State.Columns {
				columns = append(columns, Column{Title: col.ShortTitle, Width: col.Width})
			}
		case IndexAttributesGroupPodSpec:
			for _, col := range IndexColumnGroups.PodSpec.Columns {
				columns = append(columns, Column{Title: col.ShortTitle, Width: col.Width})
			}
		case IndexAttributesGroupServerlessSpec:
			for _, col := range IndexColumnGroups.ServerlessSpec.Columns {
				columns = append(columns, Column{Title: col.ShortTitle, Width: col.Width})
			}
		case IndexAttributesGroupInference:
			for _, col := range IndexColumnGroups.Inference.Columns {
				columns = append(columns, Column{Title: col.ShortTitle, Width: col.Width})
			}
		case IndexAttributesGroupOther:
			for _, col := range IndexColumnGroups.Other.Columns {
				columns = append(columns, Column{Title: col.ShortTitle, Width: col.Width})
			}
		}
	}
	return columns
}

// ExtractEssentialValues extracts essential values from an index
func ExtractEssentialValues(idx *pinecone.Index) []string {
	// Determine spec
	var spec string
	if idx.Spec.Serverless == nil {
		spec = "pod"
	} else {
		spec = "serverless"
	}

	// Determine type (for serverless indexes)
	var indexType string
	if idx.VectorType != "" {
		indexType = string(idx.VectorType)
	} else {
		indexType = "dense" // Default for pod indexes
	}

	// Get dimension
	dimension := ""
	if idx.Dimension != nil && *idx.Dimension > 0 {
		dimension = fmt.Sprintf("%d", *idx.Dimension)
	}

	return []string{
		idx.Name,
		spec,
		indexType,
		string(idx.Metric),
		dimension,
	}
}

// ExtractStateValues extracts state-related values from an index
func ExtractStateValues(idx *pinecone.Index) []string {
	// Check if protected
	protected := "no"
	if idx.DeletionProtection == pinecone.DeletionProtectionEnabled {
		protected = "yes"
	}

	status := ""
	if idx.Status != nil {
		status = string(idx.Status.State)
	}

	return []string{
		status,
		idx.Host,
		protected,
	}
}

// ExtractPodSpecValues extracts pod specification values from an index
func ExtractPodSpecValues(idx *pinecone.Index) []string {
	if idx.Spec.Pod == nil {
		return []string{"", "", "", "", ""}
	}

	return []string{
		idx.Spec.Pod.Environment,
		idx.Spec.Pod.PodType,
		fmt.Sprintf("%d", idx.Spec.Pod.Replicas),
		fmt.Sprintf("%d", idx.Spec.Pod.ShardCount),
		fmt.Sprintf("%d", idx.Spec.Pod.PodCount),
	}
}

// ExtractServerlessSpecValues extracts serverless specification values from an index
func ExtractServerlessSpecValues(idx *pinecone.Index) []string {
	if idx.Spec.Serverless == nil {
		return []string{"", ""}
	}

	return []string{
		string(idx.Spec.Serverless.Cloud),
		idx.Spec.Serverless.Region,
	}
}

// ExtractInferenceValues extracts inference-related values from an index
func ExtractInferenceValues(idx *pinecone.Index) []string {
	if idx.Embed == nil {
		return []string{"", "", "", "", ""}
	}

	embedDim := ""
	if idx.Embed.Dimension != nil && *idx.Embed.Dimension > 0 {
		embedDim = fmt.Sprintf("%d", *idx.Embed.Dimension)
	}

	// Format field map
	fieldMapStr := ""
	if idx.Embed.FieldMap != nil && len(*idx.Embed.FieldMap) > 0 {
		var fieldMapPairs []string
		for k, v := range *idx.Embed.FieldMap {
			fieldMapPairs = append(fieldMapPairs, fmt.Sprintf("%s=%v", k, v))
		}
		fieldMapStr = strings.Join(fieldMapPairs, ", ")
	}

	// Format read parameters
	readParamsStr := ""
	if idx.Embed.ReadParameters != nil && len(*idx.Embed.ReadParameters) > 0 {
		var readParamsPairs []string
		for k, v := range *idx.Embed.ReadParameters {
			readParamsPairs = append(readParamsPairs, fmt.Sprintf("%s=%v", k, v))
		}
		readParamsStr = strings.Join(readParamsPairs, ", ")
	}

	// Format write parameters
	writeParamsStr := ""
	if idx.Embed.WriteParameters != nil && len(*idx.Embed.WriteParameters) > 0 {
		var writeParamsPairs []string
		for k, v := range *idx.Embed.WriteParameters {
			writeParamsPairs = append(writeParamsPairs, fmt.Sprintf("%s=%v", k, v))
		}
		writeParamsStr = strings.Join(writeParamsPairs, ", ")
	}

	return []string{
		idx.Embed.Model,
		embedDim,
		fieldMapStr,
		readParamsStr,
		writeParamsStr,
	}
}

// ExtractOtherValues extracts other values from an index (tags, custom fields, etc.)
func ExtractOtherValues(idx *pinecone.Index) []string {
	if idx.Tags == nil || len(*idx.Tags) == 0 {
		return []string{""}
	}

	// Convert tags to a string representation showing key-value pairs
	var tagStrings []string
	for key, value := range *idx.Tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
	}
	return []string{fmt.Sprint(strings.Join(tagStrings, ", "))}
}

// ExtractValuesForIndexAttributesGroups extracts values for the specified index attribute groups from an index
func ExtractValuesForIndexAttributesGroups(idx *pinecone.Index, groups []IndexAttributesGroup) []string {
	var values []string
	for _, group := range groups {
		switch group {
		case IndexAttributesGroupEssential:
			values = append(values, ExtractEssentialValues(idx)...)
		case IndexAttributesGroupState:
			values = append(values, ExtractStateValues(idx)...)
		case IndexAttributesGroupPodSpec:
			values = append(values, ExtractPodSpecValues(idx)...)
		case IndexAttributesGroupServerlessSpec:
			values = append(values, ExtractServerlessSpecValues(idx)...)
		case IndexAttributesGroupInference:
			values = append(values, ExtractInferenceValues(idx)...)
		case IndexAttributesGroupOther:
			values = append(values, ExtractOtherValues(idx)...)
		}
	}
	return values
}

// getColumnsWithNamesForIndexAttributesGroup returns columns with both short and full names for a specific index attribute group
func getColumnsWithNamesForIndexAttributesGroup(group IndexAttributesGroup) []IndexColumn {
	switch group {
	case IndexAttributesGroupEssential:
		return IndexColumnGroups.Essential.Columns
	case IndexAttributesGroupState:
		return IndexColumnGroups.State.Columns
	case IndexAttributesGroupPodSpec:
		return IndexColumnGroups.PodSpec.Columns
	case IndexAttributesGroupServerlessSpec:
		return IndexColumnGroups.ServerlessSpec.Columns
	case IndexAttributesGroupInference:
		return IndexColumnGroups.Inference.Columns
	case IndexAttributesGroupOther:
		return IndexColumnGroups.Other.Columns
	default:
		return []IndexColumn{}
	}
}

// getValuesForIndexAttributesGroup returns values for a specific index attribute group
func getValuesForIndexAttributesGroup(idx *pinecone.Index, group IndexAttributesGroup) []string {
	switch group {
	case IndexAttributesGroupEssential:
		return ExtractEssentialValues(idx)
	case IndexAttributesGroupState:
		return ExtractStateValues(idx)
	case IndexAttributesGroupPodSpec:
		return ExtractPodSpecValues(idx)
	case IndexAttributesGroupServerlessSpec:
		return ExtractServerlessSpecValues(idx)
	case IndexAttributesGroupInference:
		return ExtractInferenceValues(idx)
	case IndexAttributesGroupOther:
		return ExtractOtherValues(idx)
	default:
		return []string{}
	}
}

// hasNonEmptyValues checks if a group has any meaningful (non-empty) values
func hasNonEmptyValues(values []string) bool {
	for _, value := range values {
		if value != "" && value != "nil" {
			return true
		}
	}
	return false
}

// filterNonEmptyIndexAttributesGroups filters out index attribute groups that have no meaningful data across all indexes
func filterNonEmptyIndexAttributesGroups(indexes []*pinecone.Index, groups []IndexAttributesGroup) []IndexAttributesGroup {
	var nonEmptyGroups []IndexAttributesGroup

	for _, group := range groups {
		hasData := false
		for _, idx := range indexes {
			values := getValuesForIndexAttributesGroup(idx, group)
			if hasNonEmptyValues(values) {
				hasData = true
				break
			}
		}
		if hasData {
			nonEmptyGroups = append(nonEmptyGroups, group)
		}
	}

	return nonEmptyGroups
}

// filterNonEmptyIndexAttributesGroupsForIndex filters out index attribute groups that have no meaningful data for a specific index
func filterNonEmptyIndexAttributesGroupsForIndex(idx *pinecone.Index, groups []IndexAttributesGroup) []IndexAttributesGroup {
	var nonEmptyGroups []IndexAttributesGroup

	for _, group := range groups {
		values := getValuesForIndexAttributesGroup(idx, group)
		if hasNonEmptyValues(values) {
			nonEmptyGroups = append(nonEmptyGroups, group)
		}
	}

	return nonEmptyGroups
}
