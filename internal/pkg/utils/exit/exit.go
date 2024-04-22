package exit

import (
	"os"
	"fmt"
)

func Error(reason error) {
	fmt.Println(reason)
	os.Exit(1)
}

func Success() {
	os.Exit(0)
}