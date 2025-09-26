package exit

import (
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
)

func Error(reason error) {
	log.Error().Msg(reason.Error())
	os.Exit(1)
}

func ErrorMsg(msg string) {
	log.Error().Msg(msg)
	os.Exit(1)
}

func Success() {
	log.Info().Msg("Exiting successfully")
	os.Exit(0)
}

func SuccessMsg(msg string) {
	log.Info().Msg(msg)
	os.Exit(0)
}
