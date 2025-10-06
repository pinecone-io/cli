package help

import (
	"regexp"
	"strings"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
)

var (
	ansiRegex = regexp.MustCompile("\x1b\\[[0-9;]*m")
)

func stripANSI(s string) string { return ansiRegex.ReplaceAllString(s, "") }

// setColor allows disabling color output for tests
// it handles restoring configuration to the previous value
func setColor(t *testing.T, on bool) {
	prev := config.Color.Get()
	config.Color.Set(on)
	t.Cleanup(func() {
		config.Color.Set(prev)
	})
}

func TestExamples_Empty(t *testing.T) {
	setColor(t, false)

	input := []string{"", "   ", "\n \t \n"}
	for _, in := range input {
		if got := Examples(in); got != "" {
			t.Fatalf("Mismatch: want empty for %q, got %q", in, got)
		}
	}
}

func TestExamples_CommandsAndComments(t *testing.T) {
	setColor(t, false)

	in := `
	    # Configure service account credentials
		pc auth configure --client-id <client-id> --client-secret <client-secret>

		# Clear configured credentials
		pc auth clear --service-account
	`
	got := Examples(in)

	want := strings.Join([]string{
		"  # Configure service account credentials",
		"  $ pc auth configure --client-id <client-id> --client-secret <client-secret>",
		"",
		"  # Clear configured credentials",
		"  $ pc auth clear --service-account",
	}, "\n")

	if got != want {
		t.Fatalf("Mismatch: want %q, got %q", want, got)
	}
}

func TestExamples_TrimsRightSpacesTabsCRLF(t *testing.T) {
	setColor(t, false)

	in := "pc one   \t\r\npc two\t \r\n"
	got := Examples(in)
	want := strings.Join([]string{
		"  $ pc one",
		"  $ pc two",
	}, "\n")

	if got != want {
		t.Fatalf("Mismatch: want %q, got %q", want, got)
	}
}

func TestExamples_PreservesInteriorBlankLines(t *testing.T) {
	setColor(t, false)

	in := "pc a\n\npc b\n"
	got := Examples(in)
	want := strings.Join([]string{
		"  $ pc a",
		"",
		"  $ pc b",
	}, "\n")

	if got != want {
		t.Fatalf("Mismatch: want %q, got %q", want, got)
	}
}
