package pcio

import (
	"fmt"
	"io"
)

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

func Printf(format string, a ...interface{}) {
	if !quiet {
		fmt.Printf(format, a...)
		return
	}
}

func Fprintf(w io.Writer, format string, a ...interface{}) {
	if !quiet {
		fmt.Fprintf(w, format, a...)
		return
	}
}

func Fprintln(w io.Writer, a ...interface{}) {
	if !quiet {
		fmt.Fprintln(w, a...)
		return
	}
}

func Fprint(w io.Writer, a ...interface{}) {
	if !quiet {
		fmt.Fprint(w, a...)
		return
	}
}

// alias Sprintf to fmt.Sprintf
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// alias Errorf to fmt.Errorf
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// alias Error to fmt.Errorf
func Error(a ...interface{}) error {
	return fmt.Errorf(fmt.Sprint(a...))
}
