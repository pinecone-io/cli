package msg

import (
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func FailMsg(format string, a ...any) {
	formatted := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, style.FailMsg(formatted))
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
