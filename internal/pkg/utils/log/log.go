package log

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	zl "github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Disabled)
	if os.Getenv("PINECONE_LOG_LEVEL") == "INFO" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	if os.Getenv("PINECONE_LOG_LEVEL") == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if os.Getenv("PINECONE_LOG_LEVEL") == "TRACE" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zl.Logger = zl.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func Logger() *zerolog.Logger {
	return &zl.Logger
}

func Trace() *zerolog.Event {
	return Logger().Trace()
}

func Debug() *zerolog.Event {
	return Logger().Debug()
}

func Info() *zerolog.Event {
	return Logger().Info()
}

func Error() *zerolog.Event {
	return Logger().Error()
}
