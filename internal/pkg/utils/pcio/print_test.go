package pcio

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = orig

	done := make(chan string)
	go func() {
		var b bytes.Buffer
		b.ReadFrom(r)
		done <- b.String()
	}()
	return <-done
}

func TestPrintJSON_EmitsWhenQuiet(t *testing.T) {
	SetQuiet(true)
	defer SetQuiet(false)

	out := captureStdout(t, func() {
		PrintJSON("{}")
	})

	assert.Equal(t, "{}\n", out)
}

func TestPrintJSON_EmitsWhenNotQuiet(t *testing.T) {
	SetQuiet(false)

	out := captureStdout(t, func() {
		PrintJSON("{}")
	})

	assert.Equal(t, "{}\n", out)
}

func TestPrintln_SuppressedWhenQuiet(t *testing.T) {
	SetQuiet(true)
	defer SetQuiet(false)

	out := captureStdout(t, func() {
		Println("should not appear")
	})

	assert.Empty(t, out)
}
