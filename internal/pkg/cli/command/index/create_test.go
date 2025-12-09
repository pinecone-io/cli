package index

import (
	"context"
	"strings"
	"testing"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

type mockIndexService struct {
	lastServerless *pinecone.CreateServerlessIndexRequest
	lastPod        *pinecone.CreatePodIndexRequest
	lastIntegrated *pinecone.CreateIndexForModelRequest
	lastBYOC       *pinecone.CreateBYOCIndexRequest
	result         *pinecone.Index
	err            error
}

func (m *mockIndexService) CreateServerlessIndex(ctx context.Context, req *pinecone.CreateServerlessIndexRequest) (*pinecone.Index, error) {
	m.lastServerless = req
	return m.result, m.err
}

func (m *mockIndexService) CreatePodIndex(ctx context.Context, req *pinecone.CreatePodIndexRequest) (*pinecone.Index, error) {
	m.lastPod = req
	return m.result, m.err
}

func (m *mockIndexService) CreateIndexForModel(ctx context.Context, req *pinecone.CreateIndexForModelRequest) (*pinecone.Index, error) {
	m.lastIntegrated = req
	return m.result, m.err
}

func (m *mockIndexService) CreateBYOCIndex(ctx context.Context, req *pinecone.CreateBYOCIndexRequest) (*pinecone.Index, error) {
	m.lastBYOC = req
	return m.result, m.err
}

func Test_runCreateIndexWithService_Serverless_Args(t *testing.T) {
	cmd := NewCreateIndexCmd()
	svc := &mockIndexService{result: &pinecone.Index{Name: "my-index"}}
	options := createIndexOptions{
		name:               "my-index",
		cloud:              "aws",
		region:             "us-east-1",
		vectorType:         "dense",
		dimension:          1536,
		metric:             "cosine",
		deletionProtection: "enabled",
		tags:               map[string]string{"tag1": "value1", "tag2": "value2"},
		sourceCollection:   "my-collection",
	}

	_, err := runCreateIndexWithService(context.Background(), cmd, svc, options)
	assert.NoError(t, err)
	assert.Nil(t, svc.lastPod)
	assert.Nil(t, svc.lastIntegrated)

	assert.Equal(t, options.name, svc.lastServerless.Name)
	assert.Equal(t, pinecone.Cloud(options.cloud), svc.lastServerless.Cloud)
	assert.Equal(t, options.region, svc.lastServerless.Region)
	assert.Equal(t, pinecone.IndexMetric(options.metric), *svc.lastServerless.Metric)
	assert.Equal(t, pinecone.DeletionProtection(options.deletionProtection), *svc.lastServerless.DeletionProtection)
	assert.Equal(t, pinecone.IndexTags(options.tags), *svc.lastServerless.Tags)
	assert.Equal(t, options.sourceCollection, *svc.lastServerless.SourceCollection)
	assert.Equal(t, options.dimension, *svc.lastServerless.Dimension)
	assert.Equal(t, options.vectorType, *svc.lastServerless.VectorType)
}

func Test_runCreateIndexWithService_Pod_Args(t *testing.T) {
	cmd := NewCreateIndexCmd()
	svc := &mockIndexService{result: &pinecone.Index{Name: "my-index"}}
	options := createIndexOptions{
		name:               "my-index",
		dimension:          1536,
		environment:        "us-east-1-aws",
		podType:            "p1.x1",
		shards:             2,
		replicas:           2,
		metric:             "cosine",
		deletionProtection: "enabled",
		tags:               map[string]string{"tag1": "value1", "tag2": "value2"},
		sourceCollection:   "my-collection",
		metadataConfig:     []string{"field1", "field2"},
	}

	_, err := runCreateIndexWithService(context.Background(), cmd, svc, options)
	assert.NoError(t, err)
	assert.Nil(t, svc.lastServerless)
	assert.Nil(t, svc.lastIntegrated)

	assert.Equal(t, options.name, svc.lastPod.Name)
	assert.Equal(t, options.dimension, svc.lastPod.Dimension)
	assert.Equal(t, options.environment, svc.lastPod.Environment)
	assert.Equal(t, options.podType, svc.lastPod.PodType)
	assert.Equal(t, options.shards, svc.lastPod.Shards)
	assert.Equal(t, options.replicas, svc.lastPod.Replicas)
	assert.Equal(t, pinecone.IndexMetric(options.metric), *svc.lastPod.Metric)
	assert.Equal(t, pinecone.DeletionProtection(options.deletionProtection), *svc.lastPod.DeletionProtection)
	assert.Equal(t, pinecone.IndexTags(options.tags), *svc.lastPod.Tags)
	assert.Equal(t, options.sourceCollection, *svc.lastPod.SourceCollection)
	assert.Equal(t, options.metadataConfig, *svc.lastPod.MetadataConfig.Indexed)
}

func Test_runCreateIndexWithService_Integrated_Args(t *testing.T) {
	cmd := NewCreateIndexCmd()
	svc := &mockIndexService{result: &pinecone.Index{Name: "my-index"}}
	options := createIndexOptions{
		name:               "my-index",
		cloud:              "aws",
		region:             "us-east-1",
		deletionProtection: "enabled",
		model:              "multilingual-e5-large",
		fieldMap:           map[string]string{"field1": "text", "field2": "text"},
		readParameters:     map[string]string{"parameter1": "value1", "parameter2": "value2"},
		writeParameters:    map[string]string{"parameter3": "value3", "parameter4": "value4"},
		tags:               map[string]string{"tag1": "value1", "tag2": "value2"},
	}

	_, err := runCreateIndexWithService(context.Background(), cmd, svc, options)
	assert.NoError(t, err)
	assert.Nil(t, svc.lastServerless)
	assert.Nil(t, svc.lastPod)

	assert.Equal(t, options.name, svc.lastIntegrated.Name)
	assert.Equal(t, pinecone.Cloud(options.cloud), svc.lastIntegrated.Cloud)
	assert.Equal(t, options.region, svc.lastIntegrated.Region)
	assert.Equal(t, pinecone.DeletionProtection(options.deletionProtection), *svc.lastIntegrated.DeletionProtection)
	assert.Equal(t, options.model, svc.lastIntegrated.Embed.Model)
	assert.Equal(t, toInterfaceMap(options.fieldMap), svc.lastIntegrated.Embed.FieldMap)
	assert.Equal(t, toInterfaceMap(options.readParameters), *svc.lastIntegrated.Embed.ReadParameters)
	assert.Equal(t, toInterfaceMap(options.writeParameters), *svc.lastIntegrated.Embed.WriteParameters)
	assert.Equal(t, pinecone.IndexTags(options.tags), *svc.lastIntegrated.Tags)
}

func TestCreateIndexOptions_DeriveIndexType(t *testing.T) {
	tests := []struct {
		name        string
		options     createIndexOptions
		expected    indexType
		expectError bool
	}{
		{
			name: "serverless - cloud, region",
			options: createIndexOptions{
				cloud:  "aws",
				region: "us-east-1",
			},
			expected: indexTypeServerless,
		},
		{
			name: "integrated - cloud, region, model",
			options: createIndexOptions{
				cloud:  "aws",
				region: "us-east-1",
				model:  "multilingual-e5-large",
			},
			expected: indexTypeIntegrated,
		},
		{
			name: "pods - environment",
			options: createIndexOptions{
				environment: "us-east-1-gcp",
			},
			expected: indexTypePod,
		},
		{
			name: "serverless - cloud and region prioritized over environment",
			options: createIndexOptions{
				cloud:       "aws",
				region:      "us-east-1",
				environment: "us-east-1-gcp",
			},
			expected: indexTypeServerless,
		},
		{
			name:        "error - no input",
			options:     createIndexOptions{},
			expectError: true,
		},
		{
			name: "error - cloud and model only",
			options: createIndexOptions{
				cloud: "aws",
				model: "multilingual-e5-large",
			},
			expectError: true,
		},
		{
			name: "error - cloud only",
			options: createIndexOptions{
				cloud: "aws",
			},
			expectError: true,
		},
		{
			name: "error - model only",
			options: createIndexOptions{
				model: "multilingual-e5-large",
			},
			expectError: true,
		},
		{
			name: "error - region only",
			options: createIndexOptions{
				region: "us-east-1",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.options.deriveIndexType()
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, got)
				}
			}
		})
	}
}

func TestCreateIndexOptions_Validate(t *testing.T) {
	tests := []struct {
		name        string
		options     createIndexOptions
		expectError bool
		errorSubstr string
	}{
		{
			name: "serverless index with name and cloud, region",
			options: createIndexOptions{
				name:  "my-index",
				cloud: "aws",
			},
			expectError: false,
		},
		{
			name: "valid - integrated index with name and cloud, region, model",
			options: createIndexOptions{
				name:   "my-index",
				cloud:  "aws",
				region: "us-east-1",
				model:  "multilingual-e5-large",
			},
		},
		{
			name: "valid - pod index with name and environment",
			options: createIndexOptions{
				name:        "my-index",
				environment: "us-east-1-gcp",
			},
			expectError: false,
		},
		{
			name:        "error - missing name",
			options:     createIndexOptions{},
			expectError: true,
			errorSubstr: "name is required",
		},
		{
			name: "error - name, cloud, region, environment all provided",
			options: createIndexOptions{
				name:        "my-index",
				cloud:       "aws",
				region:      "us-east-1",
				environment: "us-east-1-gcp",
			},
			expectError: true,
			errorSubstr: "cloud, region, and environment cannot be provided together",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errorSubstr != "" && !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("expected error to contain %q, got %q", tt.errorSubstr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
