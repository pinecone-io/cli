package index

import (
	"strings"
	"testing"
)

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
