package msg

import (
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

func FailMsg(format string, a ...any) {
	formatted := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, style.FailMsg(formatted))
}

// FailJSON emits a structured {"error": "..."} JSON object to stdout when
// jsonFlag is true, then writes the human-readable error to stderr via
// FailMsg regardless. This ensures agents capturing stdout in --json mode
// receive a machine-readable error, while human users still see styled output.
func FailJSON(jsonFlag bool, format string, a ...any) {
	if jsonFlag {
		message := fmt.Sprintf(format, a...)
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
