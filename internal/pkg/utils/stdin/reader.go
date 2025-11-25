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

// ClaimOnce reserves stdin for a single reader without reading it.
// Subsequent attempts to claim or read will error.
func ClaimOnce() error {
	if !consumed.CompareAndSwap(false, true) {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ReaderOnce returns a reader for stdin with an optional size limit.
// It marks stdin as consumed so only one caller may read.
func ReaderOnce(limit int64) (io.ReadCloser, error) {
	if err := ClaimOnce(); err != nil {
		return nil, err
	}
	var r io.Reader = os.Stdin
	if limit > 0 {
		r = io.LimitReader(os.Stdin, limit)
	}
	return io.NopCloser(r), nil
}
