package presenters

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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

// IndexAttributesGroupsToStrings converts a slice of IndexAttributesGroup to strings
func IndexAttributesGroupsToStrings(groups []IndexAttributesGroup) []string {
	strings := make([]string, len(groups))
	for i, group := range groups {
		strings[i] = string(group)
	}
	return strings
}

// StringsToIndexAttributesGroups converts a slice of strings to IndexAttributesGroup (validates input)
func StringsToIndexAttributesGroups(groups []string) []IndexAttributesGroup {
	indexGroups := make([]IndexAttributesGroup, 0, len(groups))
	validGroups := map[string]IndexAttributesGroup{
		"essential":       IndexAttributesGroupEssential,
		"state":           IndexAttributesGroupState,
		"pod_spec":        IndexAttributesGroupPodSpec,
		"serverless_spec": IndexAttributesGroupServerlessSpec,
		"inference":       IndexAttributesGroupInference,
		"other":           IndexAttributesGroupOther,
	}

	for _, group := range groups {
		if indexGroup, exists := validGroups[group]; exists {
			indexGroups = append(indexGroups, indexGroup)
		}
	}
	return indexGroups
}

// ColumnGroup represents a group of related columns for index display
type ColumnGroup struct {
	Name    string
	Columns []Column
}

// ColumnWithNames represents a table column with both short and full names
type ColumnWithNames struct {
	ShortTitle string
	FullTitle  string
	Width      int
}

// ColumnGroupWithNames represents a group of columns with both short and full names
type ColumnGroupWithNames struct {
	Name    string
	Columns []ColumnWithNames
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
		Columns: []Column{
			{Title: "NAME", Width: 20},
			{Title: "SPEC", Width: 12},
			{Title: "TYPE", Width: 8},
			{Title: "METRIC", Width: 8},
			{Title: "DIM", Width: 8},
		},
	},
	State: ColumnGroup{
		Name: "state",
		Columns: []Column{
			{Title: "STATUS", Width: 10},
			{Title: "HOST", Width: 60},
			{Title: "PROT", Width: 8},
		},
	},
	PodSpec: ColumnGroup{
		Name: "pod_spec",
		Columns: []Column{
			{Title: "ENV", Width: 12},
			{Title: "POD_TYPE", Width: 12},
			{Title: "REPLICAS", Width: 8},
			{Title: "SHARDS", Width: 8},
			{Title: "PODS", Width: 8},
		},
	},
	ServerlessSpec: ColumnGroup{
		Name: "serverless_spec",
		Columns: []Column{
			{Title: "CLOUD", Width: 12},
			{Title: "REGION", Width: 15},
		},
	},
	Inference: ColumnGroup{
		Name: "inference",
		Columns: []Column{
			{Title: "MODEL", Width: 25},
			{Title: "EMBED_DIM", Width: 10},
		},
	},
	Other: ColumnGroup{
		Name: "other",
		Columns: []Column{
			{Title: "TAGS", Width: 30},
		},
	},
}

// IndexColumnGroupsWithNames defines the available column groups with both short and full names
var IndexColumnGroupsWithNames = struct {
	Essential      ColumnGroupWithNames // Basic index information (name, spec, type, metric, dimension)
	State          ColumnGroupWithNames // Runtime state information (status, host, protection)
	PodSpec        ColumnGroupWithNames // Pod-specific configuration (environment, pod type, replicas, etc.)
	ServerlessSpec ColumnGroupWithNames // Serverless-specific configuration (cloud, region)
	Inference      ColumnGroupWithNames // Inference/embedding model information
	Other          ColumnGroupWithNames // Other information (tags, custom fields, etc.)
}{
	Essential: ColumnGroupWithNames{
		Name: "essential",
		Columns: []ColumnWithNames{
			{ShortTitle: "NAME", FullTitle: "Name", Width: 20},
			{ShortTitle: "SPEC", FullTitle: "Specification", Width: 12},
			{ShortTitle: "TYPE", FullTitle: "Vector Type", Width: 8},
			{ShortTitle: "METRIC", FullTitle: "Metric", Width: 8},
			{ShortTitle: "DIM", FullTitle: "Dimension", Width: 8},
		},
	},
	State: ColumnGroupWithNames{
		Name: "state",
		Columns: []ColumnWithNames{
			{ShortTitle: "STATUS", FullTitle: "Status", Width: 10},
			{ShortTitle: "HOST", FullTitle: "Host URL", Width: 60},
			{ShortTitle: "PROT", FullTitle: "Deletion Protection", Width: 8},
		},
	},
	PodSpec: ColumnGroupWithNames{
		Name: "pod_spec",
		Columns: []ColumnWithNames{
			{ShortTitle: "ENV", FullTitle: "Environment", Width: 12},
			{ShortTitle: "POD_TYPE", FullTitle: "Pod Type", Width: 12},
			{ShortTitle: "REPLICAS", FullTitle: "Replicas", Width: 8},
			{ShortTitle: "SHARDS", FullTitle: "Shard Count", Width: 8},
			{ShortTitle: "PODS", FullTitle: "Pod Count", Width: 8},
		},
	},
	ServerlessSpec: ColumnGroupWithNames{
		Name: "serverless_spec",
		Columns: []ColumnWithNames{
			{ShortTitle: "CLOUD", FullTitle: "Cloud Provider", Width: 12},
			{ShortTitle: "REGION", FullTitle: "Region", Width: 15},
		},
	},
	Inference: ColumnGroupWithNames{
		Name: "inference",
		Columns: []ColumnWithNames{
			{ShortTitle: "MODEL", FullTitle: "Model", Width: 25},
			{ShortTitle: "EMBED_DIM", FullTitle: "Embedding Dimension", Width: 10},
		},
	},
	Other: ColumnGroupWithNames{
		Name: "other",
		Columns: []ColumnWithNames{
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
			columns = append(columns, IndexColumnGroups.Essential.Columns...)
		case IndexAttributesGroupState:
			columns = append(columns, IndexColumnGroups.State.Columns...)
		case IndexAttributesGroupPodSpec:
			columns = append(columns, IndexColumnGroups.PodSpec.Columns...)
		case IndexAttributesGroupServerlessSpec:
			columns = append(columns, IndexColumnGroups.ServerlessSpec.Columns...)
		case IndexAttributesGroupInference:
			columns = append(columns, IndexColumnGroups.Inference.Columns...)
		case IndexAttributesGroupOther:
			columns = append(columns, IndexColumnGroups.Other.Columns...)
		}
	}
	return columns
}

// GetColumnsForIndexAttributesGroupsWithNames returns columns for the specified index attribute groups with both short and full names
func GetColumnsForIndexAttributesGroupsWithNames(groups []IndexAttributesGroup) []ColumnWithNames {
	var columns []ColumnWithNames
	for _, group := range groups {
		switch group {
		case IndexAttributesGroupEssential:
			columns = append(columns, IndexColumnGroupsWithNames.Essential.Columns...)
		case IndexAttributesGroupState:
			columns = append(columns, IndexColumnGroupsWithNames.State.Columns...)
		case IndexAttributesGroupPodSpec:
			columns = append(columns, IndexColumnGroupsWithNames.PodSpec.Columns...)
		case IndexAttributesGroupServerlessSpec:
			columns = append(columns, IndexColumnGroupsWithNames.ServerlessSpec.Columns...)
		case IndexAttributesGroupInference:
			columns = append(columns, IndexColumnGroupsWithNames.Inference.Columns...)
		case IndexAttributesGroupOther:
			columns = append(columns, IndexColumnGroupsWithNames.Other.Columns...)
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
	dimension := "nil"
	if idx.Dimension != nil && *idx.Dimension > 0 {
		dimension = pcio.Sprintf("%d", *idx.Dimension)
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

	return []string{
		string(idx.Status.State),
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
		pcio.Sprintf("%d", idx.Spec.Pod.Replicas),
		pcio.Sprintf("%d", idx.Spec.Pod.ShardCount),
		pcio.Sprintf("%d", idx.Spec.Pod.PodCount),
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
		return []string{"", ""}
	}

	embedDim := "nil"
	if idx.Embed.Dimension != nil && *idx.Embed.Dimension > 0 {
		embedDim = pcio.Sprintf("%d", *idx.Embed.Dimension)
	}

	return []string{
		idx.Embed.Model,
		embedDim,
	}
}

// ExtractOtherValues extracts other values from an index (tags, custom fields, etc.)
func ExtractOtherValues(idx *pinecone.Index) []string {
	if idx.Tags == nil || len(*idx.Tags) == 0 {
		return []string{""}
	}

	// Convert tags to a simple string representation
	// For now, just show the count, could be enhanced to show key-value pairs
	return []string{pcio.Sprintf("%d tags", len(*idx.Tags))}
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

// GetGroupDescription returns a description of what each group contains
func GetGroupDescription(group IndexAttributesGroup) string {
	switch group {
	case IndexAttributesGroupEssential:
		return "Basic index information (name, spec type, vector type, metric, dimension)"
	case IndexAttributesGroupState:
		return "Runtime state information (status, host URL, deletion protection)"
	case IndexAttributesGroupPodSpec:
		return "Pod-specific configuration (environment, pod type, replicas, shards, pod count)"
	case IndexAttributesGroupServerlessSpec:
		return "Serverless-specific configuration (cloud provider, region)"
	case IndexAttributesGroupInference:
		return "Inference/embedding model information (model name, embedding dimension)"
	case IndexAttributesGroupOther:
		return "Other information (tags, custom fields, etc.)"
	default:
		return ""
	}
}

// getColumnsWithNamesForIndexAttributesGroup returns columns with both short and full names for a specific index attribute group
func getColumnsWithNamesForIndexAttributesGroup(group IndexAttributesGroup) []ColumnWithNames {
	switch group {
	case IndexAttributesGroupEssential:
		return IndexColumnGroupsWithNames.Essential.Columns
	case IndexAttributesGroupState:
		return IndexColumnGroupsWithNames.State.Columns
	case IndexAttributesGroupPodSpec:
		return IndexColumnGroupsWithNames.PodSpec.Columns
	case IndexAttributesGroupServerlessSpec:
		return IndexColumnGroupsWithNames.ServerlessSpec.Columns
	case IndexAttributesGroupInference:
		return IndexColumnGroupsWithNames.Inference.Columns
	case IndexAttributesGroupOther:
		return IndexColumnGroupsWithNames.Other.Columns
	default:
		return []ColumnWithNames{}
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
