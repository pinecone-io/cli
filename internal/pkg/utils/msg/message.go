package msg

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

// ansiEscape matches ANSI terminal escape sequences.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// backtickCode matches backtick-wrapped code spans added by style.Code when
// colors are disabled (e.g. `pc target`), capturing the inner text.
var backtickCode = regexp.MustCompile("`([^`]*)`")

func FailMsg(format string, a ...any) {
	formatted := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, style.FailMsg(formatted))
}

// FailJSON emits a structured {"error": "..."} JSON object to stdout when
// jsonFlag is true, then writes the human-readable error to stderr via
// FailMsg regardless. This ensures agents capturing stdout in --json mode
// receive a machine-readable error, while human users still see styled output.
//
// ANSI escape sequences are stripped from the JSON error value because callers
// often pass style.Emphasis/style.Code arguments that are evaluated before this
// function runs. In a typical --json pipeline (cmd --json | jq .) stderr is
// still a TTY so colorEnabled() returns true, meaning those arguments already
// contain escape codes by the time we format the string. Stripping here keeps
// the JSON value clean regardless of terminal state.
func FailJSON(jsonFlag bool, format string, a ...any) {
	if jsonFlag {
		message := fmt.Sprintf(format, a...)
		message = ansiEscape.ReplaceAllString(message, "")
		message = backtickCode.ReplaceAllString(message, "$1")
		message = strings.TrimSpace(message)
		fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
			Error string `json:"error"`
		}{Error: message}))
	}
	FailMsg(format, a...)
}

func SuccessMsg(format string, a ...any) {
	formatted := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, style.SuccessMsg(formatted))
}

func WarnMsg(format string, a ...any) {
	formatted := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, style.WarnMsg(formatted))
}

func InfoMsg(format string, a ...any) {
	formatted := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, style.InfoMsg(formatted))
}

func HintMsg(format string, a ...any) {
	formatted := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, style.Hint(formatted))
}

func Blank() {
	fmt.Fprintln(os.Stderr, "")
}
