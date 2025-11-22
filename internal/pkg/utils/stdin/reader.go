package stdin

import (
	"io"
	"os"
	"sync/atomic"
)

var consumed atomic.Bool

// ReadAllOnce reads all data from stdin exactly once across the process.
// If stdin was already consumed by another reader, it returns an error.
func ReadAllOnce() ([]byte, error) {
	if !consumed.CompareAndSwap(false, true) {
		return nil, io.ErrUnexpectedEOF
	}
	return io.ReadAll(os.Stdin)
}

// HasPipedStdin returns true if stdin is a pipe (not a TTY).
func HasPipedStdin() bool {
	fi, _ := os.Stdin.Stat()
	return fi != nil && (fi.Mode()&os.ModeCharDevice) == 0
}
