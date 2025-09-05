package pcio

import (
	"fmt"
	"io"
)

// Package pcio provides output functions that respect the global quiet mode.
//
// USAGE GUIDELINES:
//
// Use pcio functions for:
// - User-facing messages (success, error, warning, info)
// - Progress indicators and status updates
// - Interactive prompts and confirmations
// - Help text and documentation
// - Any output that should be suppressed with -q flag
//
// Use fmt functions for:
// - Data output from informational commands (list, describe)
// - JSON output that should always be displayed
// - Table rendering and structured data display
// - String formatting (Sprintf, Errorf, Error)
// - Any output that should NOT be suppressed with -q flag
//
// The purpose of this package is to stub out the fmt package so that
// the -q quiet mode can be implemented in a consistent way across all
// commands.

var quiet bool

func SetQuiet(q bool) {
	quiet = q
}

func Println(a ...any) {
	if !quiet {
		fmt.Println(a...)
		return
	}
}

func Print(a any) {
	if !quiet {
		fmt.Print(a)
		return
	}
}

func Printf(format string, a ...any) {
	if !quiet {
		fmt.Printf(format, a...)
		return
	}
}

func Fprintf(w io.Writer, format string, a ...any) {
	if !quiet {
		fmt.Fprintf(w, format, a...)
		return
	}
}

func Fprintln(w io.Writer, a ...any) {
	if !quiet {
		fmt.Fprintln(w, a...)
		return
	}
}

func Fprint(w io.Writer, a ...any) {
	if !quiet {
		fmt.Fprint(w, a...)
		return
	}
}

// NOTE: The following three functions are aliases to `fmt` functions and do not check the quiet flag.
// This creates inconsistency with the guidelines to use `fmt` directly (not `pcio`) for non-quiet output.
// These wrappers are kept for now because:
// 1) They don't break quiet mode behavior (they're just aliases)
// 2) A mass refactoring would require updating 100+ usages across the codebase

// alias Sprintf to fmt.Sprintf
func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

// alias Errorf to fmt.Errorf
func Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

// alias Error to fmt.Errorf
func Error(a ...any) error {
	return fmt.Errorf("%s", fmt.Sprint(a...))
}
