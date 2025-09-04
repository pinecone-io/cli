package msg

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

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
func WarnMsgMultiLine(messages ...string) {
	if len(messages) == 0 {
		return
	}

	// Create a proper multi-line warning box
	formatted := style.WarnMsgMultiLine(messages...)
	pcio.Println("\n" + formatted + "\n")
}

// SuccessMsgMultiLine displays multiple success messages in a single message box
func SuccessMsgMultiLine(messages ...string) {
	if len(messages) == 0 {
		return
	}

	// Create a proper multi-line success box
	formatted := style.SuccessMsgMultiLine(messages...)
	pcio.Println("\n" + formatted + "\n")
}

// InfoMsgMultiLine displays multiple info messages in a single message box
func InfoMsgMultiLine(messages ...string) {
	if len(messages) == 0 {
		return
	}

	// Create a proper multi-line info box
	formatted := style.InfoMsgMultiLine(messages...)
	pcio.Println("\n" + formatted + "\n")
}

// FailMsgMultiLine displays multiple error messages in a single message box
func FailMsgMultiLine(messages ...string) {
	if len(messages) == 0 {
		return
	}

	// Create a proper multi-line error box
	formatted := style.FailMsgMultiLine(messages...)
	pcio.Println("\n" + formatted + "\n")
}
