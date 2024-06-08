package style

import (
	"time"

	"github.com/briandowns/spinner"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

var (
	spinnerTextEllipsis = "..."
	spinnerTextDone     = StatusGreen("done")
	spinnerTextFailed   = StatusRed("failed")

	spinnerColor = "blue"
)

func Waiting(fn func() error) error {
	return loading("", "", "", fn)
}

func Spinner(text string, fn func() error) error {
	initialMsg := text + spinnerTextEllipsis + "\n"
	doneMsg := initialMsg + spinnerTextDone + "\n"
	failMsg := initialMsg + spinnerTextFailed + "\n"

	return loading(initialMsg, doneMsg, failMsg, fn)
}

func loading(initialMsg, doneMsg, failMsg string, fn func() error) error {
	done := make(chan struct{})
	errc := make(chan error)
	go func() {
		defer close(done)

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(pcio.Messages))
		s.Prefix = initialMsg
		s.FinalMSG = doneMsg
		s.HideCursor = true
		s.Writer = pcio.Messages

		if err := s.Color(spinnerColor); err != nil {
			exit.Error(err)
		}

		s.Start()
		err := <-errc
		if err != nil {
			s.FinalMSG = failMsg
		}

		s.Stop()
	}()

	err := fn()
	errc <- err
	<-done
	return err
}
