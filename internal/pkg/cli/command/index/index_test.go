package index

import "testing"

func TestDeriveIndexType(t *testing.T) {
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
			name: "serverless - prioritized with environment",
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
