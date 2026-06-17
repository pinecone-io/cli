// Package confirm provides interactive confirmation prompts for destructive
// operations such as deletes.
package confirm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
)

// Deletion emits the provided warning lines to stderr, prompts the user to
// confirm a destructive action, and exits successfully (canceling the
// operation) if they decline. Each warning is printed verbatim, so callers may
// pre-format them (e.g. with style.Emphasis) without worrying about format
// directives.
//
// Callers are responsible for skipping this prompt when --skip-confirmation or
// --json is set.
func Deletion(warnings ...string) {
	for _, w := range warnings {
		msg.WarnMsg("%s", w)
	}

	fmt.Fprint(os.Stderr, "Do you want to continue? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		msg.FailMsg("Error reading input: %v", err)
		exit.Error(err, "Error reading input")
	}

	switch strings.TrimSpace(strings.ToLower(input)) {
	case "y", "yes":
		return
	default:
		msg.InfoMsg("Operation canceled.")
		exit.Success()
	}
}
