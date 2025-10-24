package exit

import (
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/rs/zerolog"
)

// ExitHandler interface for dependency injection / testing
type ExitHandler interface {
	Exit(code int)
}

// defaultExitHandler is the default implementation of the ExitHandler interface
// This is mocked and replaced in unit tests with setExitHandler and resetExitHandler
type defaultExitHandler struct{}

func (h *defaultExitHandler) Exit(code int) {
	os.Exit(code)
}

var exitHandler ExitHandler = &defaultExitHandler{}

func setExitHandler(handler ExitHandler) {
	exitHandler = handler
}

func resetExitHandler() {
	exitHandler = &defaultExitHandler{}
}

// exitEvent is a wrapper around zerolog.Event that adds a code and exits the program
// when Msg, Msgf, or Send is called
type exitEvent struct {
	*zerolog.Event
	code int
}

// Logs the message and exits with the exitEvent.code
func (e *exitEvent) Msg(msg string) {
	e.Event.Msg(msg)
	exitHandler.Exit(e.code)
}

// Logs the formatted message and exits with the exitEvent.code
func (e *exitEvent) Msgf(f string, v ...any) {
	e.Event.Msgf(f, v...)
	exitHandler.Exit(e.code)
}

// Equivalent to calling Msg("") then exiting with the exitEvent.code
func (e *exitEvent) Send() {
	e.Event.Send()
	exitHandler.Exit(e.code)
}

// Returns a new exitEvent/zerolog.Event with error level and code 1
func Error() *exitEvent {
	return &exitEvent{Event: log.Error(), code: 1}
}

// Returns a new exitEvent/zerolog.Event with info level and code 0
func Success() *exitEvent {
	return &exitEvent{Event: log.Info(), code: 0}
}

// Convenience function for printing a success message and exiting
func SuccessMsg(msg string) {
	Success().Msg(msg)
}

// Convenience function for printing an error message and exiting
func ErrorMsg(msg string) {
	Error().Msg(msg)
}
