package vector

import (
	"testing"
)

func TestParseUpsertBody_JSONL(t *testing.T) {
	jsonl := `{"id":"a","values":[1,2,3]}
{"id":"b","values":[4,5,6]}
`
	payload, err := parseUpsertBody([]byte(jsonl))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payload == nil || len(payload.Vectors) != 2 {
		t.Fatalf("expected 2 vectors, got %+v", payload)
	}
	if payload.Vectors[0].Id != "a" || payload.Vectors[1].Id != "b" {
		t.Fatalf("unexpected ids: %v, %v", payload.Vectors[0].Id, payload.Vectors[1].Id)
	}
}
