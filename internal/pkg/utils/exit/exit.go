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

func (e *exitEvent) WithCode(code int) *exitEvent       { e.code = code; return e }
func (e *exitEvent) Msg(msg string) *exitEvent          { e.Event.Msg(msg); return e }
func (e *exitEvent) Msgf(f string, v ...any) *exitEvent { e.Event.Msgf(f, v...); return e }
func (e *exitEvent) Send()                              { e.Event.Send(); os.Exit(e.code) }

func Error() *exitEvent {
	return &exitEvent{Event: log.Error(), code: 1}
}

func Success() *exitEvent {
	return &exitEvent{Event: log.Info(), code: 0}
}

func SuccessMsg(msg string) {
	log.Info().Msg(msg)
	os.Exit(0)
}

func ErrorMsg(msg string) {
	log.Error().Msg(msg)
	os.Exit(1)
}
