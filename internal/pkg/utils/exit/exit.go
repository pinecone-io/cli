package exit

import (
	"fmt"
	"os"
)

func Error(reason error) {
	fmt.Println(reason)
	os.Exit(1)
}

func ErrorMsg(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func Success() {
	os.Exit(0)
}
