package presenters

import (
	"testing"
)

func TestDisplayOrNone(t *testing.T) {
	str := "test"
	num := 123
	boolean := true
	emptyStr := ""
	var nilStr *string
	var nilInt *int
	var nilBool *bool
	var nilInterface any
	var nilMap map[string]string
	var nilSlice []int
	var nonNilInterface any = "wrapped"

	emptyReplacement := "<none>"

	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{name: "string", input: str, expected: str},
		{name: "empty string", input: emptyStr, expected: ""},
		{name: "nil *string", input: nilStr, expected: emptyReplacement},
		{name: "nil *int", input: nilInt, expected: emptyReplacement},
		{name: "*string", input: &str, expected: str},
		{name: "*int", input: &num, expected: "123"},
		{name: "boolean", input: boolean, expected: "true"},
		{name: "*boolean", input: &boolean, expected: "true"},
		{name: "nil *boolean", input: nilBool, expected: emptyReplacement},
		{name: "nil interface", input: nilInterface, expected: emptyReplacement},
		{name: "non-nil interface", input: nonNilInterface, expected: "wrapped"},
		{name: "nil interface pointer", input: &nilInterface, expected: emptyReplacement},
		{name: "nil map", input: nilMap, expected: "{}"},
		{name: "map value", input: map[string]string{"k": "v"}, expected: `{"k":"v"}`},
		{name: "nil map pointer", input: (*map[string]string)(nil), expected: emptyReplacement},
		{name: "nil slice", input: nilSlice, expected: "[]"},
		{name: "slice value", input: []int{1, 2}, expected: "[1,2]"},
		{name: "nil slice pointer", input: (*[]int)(nil), expected: emptyReplacement},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DisplayOrNone(tt.input)
			if result != tt.expected {
				t.Errorf("DisplayOrNone(%v) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
