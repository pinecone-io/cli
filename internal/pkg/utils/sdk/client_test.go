package sdk

import "testing"

func Test_cliSourceTag(t *testing.T) {
	t.Run("returns default when env unset", func(t *testing.T) {
		t.Setenv("PINECONE_CLI_ATTRIBUTION", "")
		if got := cliSourceTag(); got != CLISourceTag {
			t.Fatalf("expected %q, got %q", CLISourceTag, got)
		}
	})

	t.Run("appends suffix when env set", func(t *testing.T) {
		t.Setenv("PINECONE_CLI_ATTRIBUTION", "extra")
		expected := CLISourceTag + "_extra"
		if got := cliSourceTag(); got != expected {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	})
}
