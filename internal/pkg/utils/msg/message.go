package msg

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func FailMsg(format string, a ...interface{}) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println(style.FailMsg(formatted))
}

func SuccessMsg(format string, a ...interface{}) {
	formatted := pcio.Sprintf(format, a...)
	// Add a newline to the end of the message if missing
	if formatted[len(formatted)-1] != '\n' {
		formatted += "\n"
	}

	pcio.Println(style.SuccessMsg(formatted))
}
