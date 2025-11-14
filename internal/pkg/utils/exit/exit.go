package exit

import (
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
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

// Exit functions
func Success() {
	exitHandler.Exit(0)
}

func SuccessMsg(msg string) {
	log.Info().Msg(msg)
	exitHandler.Exit(0)
}

func Successf(format string, args ...any) {
	log.Info().Msgf(format, args...)
	exitHandler.Exit(0)
}

func ErrorMsg(msg string) {
	log.Error().Msg(msg)
	exitHandler.Exit(1)
}

func ErrorMsgf(format string, args ...any) {
	log.Error().Msgf(format, args...)
	exitHandler.Exit(1)
}

func Error(err error, msg string) {
	if err != nil {
		log.Error().Err(err).Msg(msg)
	} else {
		log.Error().Msg(msg)
	}
	exitHandler.Exit(1)
}

func Errorf(err error, format string, args ...any) {
	if err != nil {
		log.Error().Err(err).Msgf(format, args...)
	} else {
		log.Error().Msgf(format, args...)
	}
	exitHandler.Exit(1)
}
