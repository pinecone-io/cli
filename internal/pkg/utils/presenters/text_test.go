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
	var nonNilInterface any = "wrapped"

	emptyReplacement := "<none>"

	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "string",
			input:    str,
			expected: str,
		},
		{
			name:     "empty string",
			input:    emptyStr,
			expected: "",
		},
		{
			name:     "nil *string",
			input:    nilStr,
			expected: emptyReplacement,
		},
		{
			name:     "nil *int",
			input:    nilInt,
			expected: emptyReplacement,
		},
		{
			name:     "*string",
			input:    &str,
			expected: str,
		},
		{
			name:     "*int",
			input:    &num,
			expected: num,
		},
		{
			name:     "boolean",
			input:    boolean,
			expected: boolean,
		},
		{
			name:     "*boolean",
			input:    &boolean,
			expected: boolean,
		},
		{
			name:     "nil *boolean",
			input:    nilBool,
			expected: emptyReplacement,
		},
		{
			name:     "nil interface",
			input:    nilInterface,
			expected: emptyReplacement,
		},
		{
			name:     "non-nil interface",
			input:    nonNilInterface,
			expected: nonNilInterface,
		},
		{
			name:     "nil interface pointer",
			input:    &nilInterface,
			expected: emptyReplacement,
		},
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
