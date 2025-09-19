package msg

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

// FailMsg displays an error message to the user.
// Uses pcio functions so the message is suppressed with -q flag.
func FailMsg(format string, a ...any) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println("\n" + style.FailMsg(formatted) + "\n")
}

func SuccessMsg(format string, a ...any) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println("\n" + style.SuccessMsg(formatted) + "\n")
}

func WarnMsg(format string, a ...any) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println("\n" + style.WarnMsg(formatted) + "\n")
}

func InfoMsg(format string, a ...any) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println("\n" + style.InfoMsg(formatted) + "\n")
}

func HintMsg(format string, a ...any) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println(style.Hint(formatted))
}

// WarnMsgMultiLine displays multiple warning messages in a single message box
func FailMsgMultiLine(messages ...string) {
	if len(messages) == 0 {
		return
	}

	if len(messages) == 1 {
		FailMsg(messages[0])
		return
	}

	// Multi-line - use existing multi-line styling
	formatted := style.FailMsgMultiLine(messages...)
	pcio.Println("\n" + formatted + "\n")
}

// SuccessMsgMultiLine displays multiple success messages in a single message box
func SuccessMsgMultiLine(messages ...string) {
	if len(messages) == 0 {
		return
	}

	if len(messages) == 1 {
		SuccessMsg(messages[0])
		return
	}

	// Multi-line - use existing multi-line styling
	formatted := style.SuccessMsgMultiLine(messages...)
	pcio.Println("\n" + formatted + "\n")
}

// InfoMsgMultiLine displays multiple info messages in a single message box
func WarnMsgMultiLine(messages ...string) {
	if len(messages) == 0 {
		return
	}

	if len(messages) == 1 {
		WarnMsg(messages[0])
		return
	}

	// Multi-line - use existing multi-line styling
	formatted := style.WarnMsgMultiLine(messages...)
	pcio.Println("\n" + formatted + "\n")
}

// FailMsgMultiLine displays multiple error messages in a single message box
func InfoMsgMultiLine(messages ...string) {
	if len(messages) == 0 {
		return
	}

	if len(messages) == 1 {
		InfoMsg(messages[0])
		return
	}

	// Multi-line - use existing multi-line styling
	formatted := style.InfoMsgMultiLine(messages...)
	pcio.Println("\n" + formatted + "\n")
}
