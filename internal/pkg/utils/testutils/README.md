# Test Utilities

This package provides reusable test utilities for testing CLI commands, particularly for common patterns like the `--json` flag and argument validation.

## File Organization

- `testutils.go` - Complex testing utilities (`TestCommandArgsAndFlags`)
- `assertions.go` - Simple assertion utilities (`AssertCommandUsage`, `AssertJSONFlag`)
- `index_validation.go` - Index name validation utilities (`GetIndexNameValidationTests`)

## Generic Command Testing

The most powerful utility is `TestCommandArgsAndFlags` which provides a generic way to test any command's argument validation and flag handling:

```go
func TestMyCommand_ArgsValidation(t *testing.T) {
    cmd := NewMyCommand()

    tests := []testutils.CommandTestConfig{
        {
            Name:         "valid - single argument with flag",
            Args:         []string{"my-arg"},
            Flags:        map[string]string{"json": "true"},
            ExpectError:  false,
            ExpectedArgs: []string{"my-arg"},
            ExpectedFlags: map[string]interface{}{
                "json": true,
            },
        },
        {
            Name:        "error - no arguments",
            Args:        []string{},
            Flags:       map[string]string{},
            ExpectError: true,
            ErrorSubstr: "please provide an argument",
        },
    }

    testutils.TestCommandArgsAndFlags(t, cmd, tests)
}
```

## JSON Flag Testing

The `--json` flag is used across many commands in the CLI. The JSON utility answers one simple question: **"Does my command have a properly configured `--json` flag?"**

### JSON Flag Test

```go
func TestMyCommand_Flags(t *testing.T) {
    cmd := NewMyCommand()

    // Test that the command has a properly configured --json flag
    testutils.AssertJSONFlag(t, cmd)
}
```

This single function verifies that the `--json` flag is:

- Properly defined
- Boolean type
- Optional (not required)
- Has a description mentioning "json"
- Can be set to true/false

## Command Usage Testing

The `AssertCommandUsage` utility tests that a command has proper usage metadata:

```go
func TestMyCommand_Usage(t *testing.T) {
    cmd := NewMyCommand()

    // Test that the command has proper usage metadata
    testutils.AssertCommandUsage(t, cmd, "my-command <arg>", "domain")
}
```

This function verifies that the command has:

- Correct usage string format
- Non-empty short description
- Description mentions the expected domain

## Index Name Validation

For commands that take an index name as a positional argument (like `describe`, `delete`, etc.), there are specialized utilities:

### Index Name Validator

**Basic approach (preset tests only):**

```go
func TestMyIndexCommand_ArgsValidation(t *testing.T) {
    cmd := NewMyIndexCommand()

    // Get preset index name validation tests
    tests := testutils.GetIndexNameValidationTests()

    // Use the generic test utility
    testutils.TestCommandArgsAndFlags(t, cmd, tests)
}
```

**Advanced approach (preset + custom tests):**

```go
func TestMyIndexCommand_ArgsValidation(t *testing.T) {
    cmd := NewMyIndexCommand()

    // Get preset index name validation tests
    tests := testutils.GetIndexNameValidationTests()

    // Add custom tests for this specific command
    customTests := []testutils.CommandTestConfig{
        {
            Name:        "valid - with custom flag",
            Args:        []string{"my-index"},
            Flags:       map[string]string{"custom-flag": "value"},
            ExpectError: false,
        },
    }

    // Combine preset tests with custom tests
    allTests := append(tests, customTests...)

    // Use the generic test utility
    testutils.TestCommandArgsAndFlags(t, cmd, allTests)
}
```

### Testing Flags Separately

**For commands with --json flag:**

```go
func TestMyIndexCommand_Flags(t *testing.T) {
    cmd := NewMyIndexCommand()

    // Test the --json flag specifically
    testutils.AssertJSONFlag(t, cmd)
}
```

**For commands with custom flags:**

```go
func TestMyIndexCommand_Flags(t *testing.T) {
    cmd := NewMyIndexCommand()

    // Test custom flags using the generic utility
    tests := []testutils.CommandTestConfig{
        {
            Name:         "valid - with custom flag",
            Args:         []string{"my-index"},
            Flags:        map[string]string{"custom-flag": "value"},
            ExpectError:  false,
            ExpectedArgs: []string{"my-index"},
            ExpectedFlags: map[string]interface{}{
                "custom-flag": "value",
            },
        },
    }

    testutils.TestCommandArgsAndFlags(t, cmd, tests)
}
```

### Index Name Validator Function

```go
func NewMyIndexCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "my-command <name>",
        Short: "Description of my command",
        Args: testutils.CreateIndexNameValidator(), // Reusable validator
        Run: func(cmd *cobra.Command, args []string) {
            // Command logic here
        },
    }
    return cmd
}
```

The index name validator handles:

- Empty string validation
- Whitespace-only validation
- Multiple argument validation
- No argument validation

## Available Functions

### Generic Command Testing

- `TestCommandArgsAndFlags(t, cmd, tests)` - Generic utility to test any command's argument validation and flag handling
- `CommandTestConfig` - Configuration struct for defining test cases
- `AssertCommandUsage(t, cmd, expectedUsage, expectedDomain)` - Tests that a command has proper usage metadata

### Index Name Validation

- `GetIndexNameValidationTests()` - Returns preset test cases for index name validation

### JSON Flag Specific

- `AssertJSONFlag(t, cmd)` - Verifies that the command has a properly configured `--json` flag (definition, type, optional, description, functionality)

## Supported Flag Types

The generic utility supports all common flag types:

- `bool` - Boolean flags
- `string` - String flags
- `int`, `int64` - Integer flags
- `float64` - Float flags
- `stringSlice`, `intSlice` - Slice flags

## Usage in Commands

Any command that follows the standard cobra pattern can use these utilities. The generic utilities are particularly useful for commands with:

- Positional arguments
- Multiple flags of different types
- Custom argument validation logic

## Example

See `internal/pkg/cli/command/index/describe_test.go` for a complete example of how to use these utilities.
