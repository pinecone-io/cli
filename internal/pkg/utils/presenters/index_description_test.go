package presenters

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func TestPrintDescribeIndexTable(t *testing.T) {
	dim3 := int32(3)
	podDim := int32(4)
	sourceCollection := "sc-1"

	tests := []struct {
		name     string
		index    *pinecone.Index
		expected string
	}{
		{
			name:     "nil index prints empty state",
			index:    nil,
			expected: "No index details available.",
		},
		{
			name: "serverless spec with schema and read capacity",
			index: &pinecone.Index{
				Name:               "test-index",
				Metric:             pinecone.Cosine,
				VectorType:         "dense",
				Dimension:          &dim3,
				DeletionProtection: pinecone.DeletionProtectionEnabled,
				Status: &pinecone.IndexStatus{
					Ready: true,
					State: pinecone.Ready,
				},
				Host: "https://example.com",
				Spec: &pinecone.IndexSpec{
					Serverless: &pinecone.ServerlessSpec{
						Cloud:  pinecone.Aws,
						Region: "us-east-1",
						Schema: &pinecone.MetadataSchema{
							Fields: map[string]pinecone.MetadataSchemaField{
								"genre": {Filterable: true},
							},
						},
						ReadCapacity: &pinecone.ReadCapacity{
							OnDemand: &pinecone.ReadCapacityOnDemand{
								Status: pinecone.ReadCapacityStatus{State: "Ready"},
							},
						},
					},
				},
			},
			expected: `
ATTRIBUTE	VALUE
Name	test-index
Dimension	3
Metric	cosine
Deletion Protection	enabled
Vector Type	dense

State	Ready
Ready	true
Host	https://example.com
Private Host	<none>

Spec	serverless
Cloud	aws
Region	us-east-1
Source Collection	<none>
Schema	{"fields":{"genre":{"filterable":true}}}
Read Capacity	{"on_demand":{"status":{"state":"Ready"}}}
`,
		},
		{
			name: "pod spec with metadata config",
			index: &pinecone.Index{
				Name:               "pod-index",
				Metric:             pinecone.Dotproduct,
				VectorType:         "sparse",
				Dimension:          &podDim,
				DeletionProtection: pinecone.DeletionProtectionDisabled,
				Host:               "api.pinecone.io",
				PrivateHost:        &sourceCollection,
				Spec: &pinecone.IndexSpec{
					Pod: &pinecone.PodSpec{
						Environment:      "gcp-starter",
						PodType:          "p1.x1",
						Replicas:         2,
						ShardCount:       1,
						PodCount:         2,
						MetadataConfig:   &pinecone.PodSpecMetadataConfig{Indexed: &[]string{"genre"}},
						SourceCollection: &sourceCollection,
					},
				},
			},
			expected: `
ATTRIBUTE	VALUE
Name	pod-index
Dimension	4
Metric	dotproduct
Deletion Protection	disabled
Vector Type	sparse

State	<none>
Ready	<none>
Host	api.pinecone.io
Private Host	sc-1

Spec	pod
Environment	gcp-starter
PodType	p1.x1
Replicas	2
ShardCount	1
PodCount	2
MetadataConfig	{"indexed":["genre"]}
Source Collection	sc-1
`,
		},
		{
			name: "spec missing with embed and tags",
			index: &pinecone.Index{
				Name:       "embed-index",
				Metric:     pinecone.Cosine,
				VectorType: "dense",
				Host:       "embed-host",
				Spec:       nil,
				Embed: &pinecone.IndexEmbed{
					Model: "text-embedding-3-small",
					FieldMap: pointerToMap(map[string]interface{}{
						"title": "text",
					}),
					ReadParameters: pointerToMap(map[string]interface{}{
						"top_k": 5,
					}),
					WriteParameters: pointerToMap(map[string]interface{}{
						"namespace": "news",
					}),
				},
				Tags: pointerToTags(pinecone.IndexTags{
					"env":  "dev",
					"team": "search",
				}),
			},
			expected: `
ATTRIBUTE	VALUE
Name	embed-index
Dimension	<none>
Metric	cosine
Deletion Protection	disabled
Vector Type	dense

State	<none>
Ready	<none>
Host	embed-host
Private Host	<none>

Spec	<none>

Model	text-embedding-3-small
Field Map	{"title":"text"}
Read Parameters	{"top_k":5}
Write Parameters	{"namespace":"news"}
Tags	{"env":"dev","team":"search"}
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeOutput(captureOutput(t, func() {
				PrintDescribeIndexTable(tt.index)
			}))
			want := strings.TrimSpace(tt.expected)

			if got != want {
				t.Fatalf("unexpected output\nwant:\n%q\ngot:\n%q", want, got)
			}
		})
	}
}

var (
	ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	gapRegexp  = regexp.MustCompile(`\s{2,}`)
)

func normalizeOutput(output string) string {
	clean := stripANSI(output)
	lines := strings.Split(strings.TrimRight(clean, "\n"), "\n")

	for i, line := range lines {
		trimmed := strings.TrimRight(line, " ")
		if strings.TrimSpace(trimmed) == "" {
			lines[i] = ""
			continue
		}

		parts := gapRegexp.Split(trimmed, 2)
		if len(parts) == 2 {
			lines[i] = strings.TrimSpace(parts[0]) + "\t" + strings.TrimSpace(parts[1])
			continue
		}

		lines[i] = trimmed
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

func pointerToMap(m map[string]interface{}) *map[string]interface{} {
	return &m
}

func pointerToTags(tags pinecone.IndexTags) *pinecone.IndexTags {
	return &tags
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	_ = r.Close()
	return string(out)
}
