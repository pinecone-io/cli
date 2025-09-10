package index

// IndexSpec represents the type of index (serverless, pod) as per what the server recognizes
type IndexSpec string

const (
	IndexSpecServerless IndexSpec = "serverless"
	IndexSpecPod        IndexSpec = "pod"
)

// IndexCreateFlow represents the type of index for the creation flow
type IndexCreateFlow int

const (
	Serverless IndexCreateFlow = iota
	Pod
	Integrated
)

type EmbeddingModel string

const (
	LlamaTextEmbedV2        EmbeddingModel = "llama-text-embed-v2"
	MultilingualE5Large     EmbeddingModel = "multilingual-e5-large"
	PineconeSparseEnglishV0 EmbeddingModel = "pinecone-sparse-english-v0"
)

// Option represents a configuration option with its value and whether it was inferred
type Option[T any] struct {
	Value    T
	Inferred bool
}

// CreateOptions represents the configuration for creating an index
type CreateOptions struct {
	Name               Option[string]
	Serverless         Option[bool]
	Pod                Option[bool]
	VectorType         Option[string]
	Cloud              Option[string]
	Region             Option[string]
	SourceCollection   Option[string]
	Environment        Option[string]
	PodType            Option[string]
	Shards             Option[int32]
	Replicas           Option[int32]
	MetadataConfig     Option[[]string]
	Model              Option[string]
	FieldMap           Option[map[string]string]
	ReadParameters     Option[map[string]string]
	WriteParameters    Option[map[string]string]
	Dimension          Option[int32]
	Metric             Option[string]
	DeletionProtection Option[string]
	Tags               Option[map[string]string]
}

// GetSpec determines the index specification type based on the flags
func (c *CreateOptions) GetSpec() IndexSpec {
	if c.Pod.Value && !c.Serverless.Value {
		return IndexSpecPod
	}

	if c.Serverless.Value && !c.Pod.Value {
		return IndexSpecServerless
	}
	return ""
}

// GetSpecString returns the spec as a string for the presenter interface
func (c *CreateOptions) GetSpecString() (string, bool) {
	spec := c.GetSpec()
	return string(spec), c.Serverless.Inferred || c.Pod.Inferred
}

func (c *CreateOptions) GetCreateFlow() IndexCreateFlow {
	if c.GetSpec() == IndexSpecPod {
		return Pod
	}

	if c.GetSpec() == IndexSpecServerless && c.Model.Value != "" {
		return Integrated
	}

	return Serverless
}

// InferredCreateOptions returns CreateOptions with inferred values applied based on the spec
func InferredCreateOptions(opts CreateOptions) CreateOptions {

	if EmbeddingModel(opts.Model.Value) == "default" || EmbeddingModel(opts.Model.Value) == "default-dense" {
		opts.Model = Option[string]{
			Value:    string(LlamaTextEmbedV2),
			Inferred: true,
		}
	}

	if EmbeddingModel(opts.Model.Value) == "default-sparse" {
		opts.Model = Option[string]{
			Value:    string(PineconeSparseEnglishV0),
			Inferred: true,
		}
	}

	if EmbeddingModel(opts.Model.Value) == LlamaTextEmbedV2 {
		opts.Serverless = Option[bool]{
			Value:    true,
			Inferred: true,
		}
		opts.FieldMap = Option[map[string]string]{
			Value:    map[string]string{"text": "text"},
			Inferred: true,
		}
		opts.ReadParameters = Option[map[string]string]{
			Value:    map[string]string{"input_type": "query", "truncate": "END"},
			Inferred: true,
		}
		opts.WriteParameters = Option[map[string]string]{
			Value:    map[string]string{"input_type": "passage", "truncate": "END"},
			Inferred: true,
		}
	}

	if EmbeddingModel(opts.Model.Value) == MultilingualE5Large {
		opts.Serverless = Option[bool]{
			Value:    true,
			Inferred: true,
		}
		opts.FieldMap = Option[map[string]string]{
			Value:    map[string]string{"text": "text"},
			Inferred: true,
		}
		opts.ReadParameters = Option[map[string]string]{
			Value:    map[string]string{"input_type": "query", "truncate": "END"},
			Inferred: true,
		}
		opts.WriteParameters = Option[map[string]string]{
			Value:    map[string]string{"input_type": "passage", "truncate": "END"},
			Inferred: true,
		}
	}

	if EmbeddingModel(opts.Model.Value) == PineconeSparseEnglishV0 {
		opts.Serverless = Option[bool]{
			Value:    true,
			Inferred: true,
		}
		opts.FieldMap = Option[map[string]string]{
			Value:    map[string]string{"text": "text"},
			Inferred: true,
		}
		opts.ReadParameters = Option[map[string]string]{
			Value:    map[string]string{"input_type": "query", "truncate": "END"},
			Inferred: true,
		}
		opts.WriteParameters = Option[map[string]string]{
			Value:    map[string]string{"input_type": "passage", "truncate": "END"},
			Inferred: true,
		}
		opts.Dimension = Option[int32]{
			Value:    0,
			Inferred: true,
		}
		opts.VectorType = Option[string]{
			Value:    "sparse",
			Inferred: true,
		}
		opts.Metric = Option[string]{
			Value:    "dotproduct",
			Inferred: true,
		}
	}

	// set serverless to true if no spec is provided
	if opts.GetSpec() == "" {
		opts.Serverless = Option[bool]{
			Value:    true,
			Inferred: true,
		}
	}

	// Set vector type to dense unless already set
	if opts.VectorType.Value == "" {
		opts.VectorType = Option[string]{
			Value:    "dense",
			Inferred: true,
		}
	}

	// set cloud to aws if serverless and no cloud is provided
	if opts.GetSpec() == IndexSpecServerless && opts.Cloud.Value == "" {
		opts.Cloud = Option[string]{
			Value:    "aws",
			Inferred: true,
		}
	}

	// Infer default region based on cloud if region is not set
	if opts.Cloud.Value != "" && opts.Region.Value == "" {
		switch opts.Cloud.Value {
		case "aws":
			opts.Region = Option[string]{
				Value:    "us-east-1",
				Inferred: true,
			}
		case "gcp":
			opts.Region = Option[string]{
				Value:    "us-central1",
				Inferred: true,
			}
		case "azure":
			opts.Region = Option[string]{
				Value:    "eastus2",
				Inferred: true,
			}
		}
	}

	if opts.GetSpec() == IndexSpecPod {
		if opts.PodType.Value == "" {
			opts.PodType = Option[string]{
				Value:    "p1.x1",
				Inferred: true,
			}
		}
		if opts.Environment.Value == "" {
			opts.Environment = Option[string]{
				Value:    "us-east-1-aws",
				Inferred: true,
			}
		}
		if opts.Shards.Value == 0 {
			opts.Shards = Option[int32]{
				Value:    1,
				Inferred: true,
			}
		}
		if opts.Replicas.Value == 0 {
			opts.Replicas = Option[int32]{
				Value:    1,
				Inferred: true,
			}
		}
	}

	if opts.VectorType.Value == "dense" && opts.Dimension.Value == 0 {
		opts.Dimension = Option[int32]{
			Value:    1024,
			Inferred: true,
		}
	}

	// metric should be dotproduct when vector type is sparse
	if opts.VectorType.Value == "sparse" && opts.Metric.Value == "" {
		opts.Metric = Option[string]{
			Value:    "dotproduct",
			Inferred: true,
		}
	}

	if opts.Metric.Value == "" {
		opts.Metric = Option[string]{
			Value:    "cosine",
			Inferred: true,
		}
	}

	if opts.DeletionProtection.Value == "" {
		opts.DeletionProtection = Option[string]{
			Value:    "disabled",
			Inferred: true,
		}
	}

	return opts
}
