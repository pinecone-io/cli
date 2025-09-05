package index

// IndexSpec represents the type of index (serverless, pod, integrated)
type IndexSpec string

const (
	IndexSpecServerless IndexSpec = "serverless"
	IndexSpecPod        IndexSpec = "pod"
)

// CreateOptions represents the configuration for creating an index
type CreateOptions struct {
	Name               string
	Serverless         bool
	Pod                bool
	VectorType         string
	Cloud              string
	Region             string
	SourceCollection   string
	Environment        string
	PodType            string
	Shards             int32
	Replicas           int32
	MetadataConfig     []string
	Model              string
	FieldMap           map[string]string
	ReadParameters     map[string]string
	WriteParameters    map[string]string
	Dimension          int32
	Metric             string
	DeletionProtection string
	Tags               map[string]string
}

// GetSpec determines the index specification type based on the flags
func (c *CreateOptions) GetSpec() IndexSpec {
	if c.Serverless && c.Pod {
		return "" // This should be caught by validation
	}
	if c.Pod {
		return IndexSpecPod
	}
	// default to serverless
	return IndexSpecServerless
}

// GetSpecString returns the spec as a string for the presenter interface
func (c *CreateOptions) GetSpecString() string {
	spec := c.GetSpec()
	return string(spec)
}
