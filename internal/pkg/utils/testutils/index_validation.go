package testutils

// GetIndexNameValidationTests returns standard index name validation test cases
// These tests focus ONLY on index name validation, without any flag assumptions
// Flags should be tested separately using the generic TestCommandArgsAndFlags utility
func GetIndexNameValidationTests() []CommandTestConfig {
	return []CommandTestConfig{
		{
			Name:        "valid - single index name",
			Args:        []string{"my-index"},
			Flags:       map[string]string{},
			ExpectError: false,
		},
		{
			Name:        "valid - index name with special characters",
			Args:        []string{"my-index-123"},
			Flags:       map[string]string{},
			ExpectError: false,
		},
		{
			Name:        "valid - index name with underscores",
			Args:        []string{"my_index_123"},
			Flags:       map[string]string{},
			ExpectError: false,
		},
		{
			Name:        "error - no arguments",
			Args:        []string{},
			Flags:       map[string]string{},
			ExpectError: true,
			ErrorSubstr: "please provide an index name",
		},
		{
			Name:        "error - multiple arguments",
			Args:        []string{"index1", "index2"},
			Flags:       map[string]string{},
			ExpectError: true,
			ErrorSubstr: "please provide only one index name",
		},
		{
			Name:        "error - three arguments",
			Args:        []string{"index1", "index2", "index3"},
			Flags:       map[string]string{},
			ExpectError: true,
			ErrorSubstr: "please provide only one index name",
		},
		{
			Name:        "error - empty string argument",
			Args:        []string{""},
			Flags:       map[string]string{},
			ExpectError: true,
			ErrorSubstr: "index name cannot be empty",
		},
		{
			Name:        "error - whitespace only argument",
			Args:        []string{"   "},
			Flags:       map[string]string{},
			ExpectError: true,
			ErrorSubstr: "index name cannot be empty",
		},
	}
}
