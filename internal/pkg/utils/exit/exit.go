package exit

import (
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/rs/zerolog"
)

type exitEvent struct {
	*zerolog.Event
	code int
}

// Logs the message and exits with the exitEvent.code
func (e *exitEvent) Msg(msg string) {
	e.Event.Msg(msg)
	os.Exit(e.code)
}

// Logs the formatted message and exits with the exitEvent.code
func (e *exitEvent) Msgf(f string, v ...any) {
	e.Event.Msgf(f, v...)
	os.Exit(e.code)
}

// Equivalent to calling Msg("") then exiting with the exitEvent.code
func (e *exitEvent) Send() {
	e.Event.Send()
	os.Exit(e.code)
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
