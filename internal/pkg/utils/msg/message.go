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
	pcio.Println(style.SuccessMsg(formatted))
}

func WarnMsg(format string, a ...interface{}) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println(style.WarnMsg(formatted))
}

func InfoMsg(format string, a ...interface{}) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println(style.InfoMsg(formatted))
}

func HintMsg(format string, a ...interface{}) {
	formatted := pcio.Sprintf(format, a...)
	pcio.Println(style.Hint(formatted))
}
