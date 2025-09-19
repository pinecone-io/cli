package testutils

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommandTestConfig represents the configuration for testing a command's arguments and flags
type CommandTestConfig struct {
	Name          string
	Args          []string
	Flags         map[string]string
	ExpectError   bool
	ErrorSubstr   string
	ExpectedArgs  []string               // Expected positional arguments after processing
	ExpectedFlags map[string]interface{} // Expected flag values after processing
}

// TestCommandArgsAndFlags provides a generic way to test any command's argument validation and flag handling
// This can be used for any command that follows the standard cobra pattern
func TestCommandArgsAndFlags(t *testing.T, cmd *cobra.Command, tests []CommandTestConfig) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Create a fresh command instance for each test
			cmdCopy := *cmd
			cmdCopy.Flags().SortFlags = false

			// Reset all flags to their default values
			cmdCopy.Flags().VisitAll(func(flag *pflag.Flag) {
				flag.Value.Set(flag.DefValue)
			})

			// Set flags if provided
			for flagName, flagValue := range tt.Flags {
				err := cmdCopy.Flags().Set(flagName, flagValue)
				if err != nil {
					t.Fatalf("failed to set flag %s=%s: %v", flagName, flagValue, err)
				}
			}

			// Test the Args validation function
			err := cmdCopy.Args(&cmdCopy, tt.Args)

			if tt.ExpectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.ErrorSubstr != "" && !strings.Contains(err.Error(), tt.ErrorSubstr) {
					t.Errorf("expected error to contain %q, got %q", tt.ErrorSubstr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				// If validation passed, test that the command would be configured correctly
				if len(tt.Args) > 0 && len(tt.ExpectedArgs) > 0 {
					// Verify positional arguments
					for i, expectedArg := range tt.ExpectedArgs {
						if i < len(tt.Args) && tt.Args[i] != expectedArg {
							t.Errorf("expected arg[%d] to be %q, got %q", i, expectedArg, tt.Args[i])
						}
					}
				}

				// Verify flag values
				for flagName, expectedValue := range tt.ExpectedFlags {
					flag := cmdCopy.Flags().Lookup(flagName)
					if flag == nil {
						t.Errorf("expected flag %s to exist", flagName)
						continue
					}

					// Get the actual flag value based on its type
					actualValue, err := getFlagValue(&cmdCopy, flagName, flag.Value.Type())
					if err != nil {
						t.Errorf("failed to get flag %s value: %v", flagName, err)
						continue
					}

					if actualValue != expectedValue {
						t.Errorf("expected flag %s to be %v, got %v", flagName, expectedValue, actualValue)
					}
				}
			}
		})
	}
}

// getFlagValue retrieves the value of a flag based on its type
func getFlagValue(cmd *cobra.Command, flagName, flagType string) (interface{}, error) {
	switch flagType {
	case "bool":
		return cmd.Flags().GetBool(flagName)
	case "string":
		return cmd.Flags().GetString(flagName)
	case "int":
		return cmd.Flags().GetInt(flagName)
	case "int64":
		return cmd.Flags().GetInt64(flagName)
	case "float64":
		return cmd.Flags().GetFloat64(flagName)
	case "stringSlice":
		return cmd.Flags().GetStringSlice(flagName)
	case "intSlice":
		return cmd.Flags().GetIntSlice(flagName)
	default:
		return cmd.Flags().GetString(flagName) // fallback to string
	}
}
