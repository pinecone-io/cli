package exit

import (
	"bytes"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type MockExitHandler struct {
	LastExitCode int
	ExitCalled   bool
	ExitCount    int
}

func (m *MockExitHandler) Exit(code int) {
	m.LastExitCode = code
	m.ExitCalled = true
	m.ExitCount++
}

func (m *MockExitHandler) Reset() {
	m.LastExitCode = 0
	m.ExitCalled = false
	m.ExitCount = 0
}

// withCapturedLogs configures zerolog to write JSON logs to a buffer and returns
// a restore function to revert logger and level after the test, let's us assert on the logged output
func withCapturedLogs(t *testing.T) (func(), *bytes.Buffer) {
	t.Helper()
	prevLogger := *log.Logger()
	prevLevel := zerolog.GlobalLevel()

	buf := &bytes.Buffer{}
	*log.Logger() = log.Logger().Output(buf)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	restore := func() {
		*log.Logger() = prevLogger
		zerolog.SetGlobalLevel(prevLevel)
	}
	return restore, buf
}

func TestSuccess_ExitsWithZero(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	Success()

	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 0, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
}

func TestSuccessMsg_LogsInfoAndExitsZero(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	restore, buf := withCapturedLogs(t)
	defer restore()

	SuccessMsg("test")

	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 0, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
	s := buf.String()
	assert.Contains(t, s, "\"level\":\"info\"")
	assert.Contains(t, s, "test")
}

func TestSuccessf_LogsInfoAndExitsZero(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	restore, buf := withCapturedLogs(t)
	defer restore()

	Successf("hello %s", "world")

	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 0, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
	s := buf.String()
	assert.Contains(t, s, "\"level\":\"info\"")
	assert.Contains(t, s, "hello world")
}

func TestErrorMsg_LogsErrorAndExitsOne(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	restore, buf := withCapturedLogs(t)
	defer restore()

	ErrorMsg("boom")

	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
	s := buf.String()
	assert.Contains(t, s, "\"level\":\"error\"")
	assert.Contains(t, s, "boom")
}

func TestErrorMsgf_LogsErrorAndExitsOne(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	restore, buf := withCapturedLogs(t)
	defer restore()

	ErrorMsgf("boom %s", "now")

	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
	s := buf.String()
	assert.Contains(t, s, "\"level\":\"error\"")
	assert.Contains(t, s, "boom now")
}

func TestError_WithErr_LogsErrorFieldAndExitsOne(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	restore, buf := withCapturedLogs(t)
	defer restore()

	Error(errors.New("some err"), "failed to do thing")

	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
	s := buf.String()
	assert.Contains(t, s, "\"level\":\"error\"")
	assert.Contains(t, s, "failed to do thing")
	assert.Contains(t, s, "some err")
}

func TestError_WithNilErr_LogsMessageAndExitsOne(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	restore, buf := withCapturedLogs(t)
	defer restore()

	Error(nil, "failed with nil err")

	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
	s := buf.String()
	assert.Contains(t, s, "\"level\":\"error\"")
	assert.Contains(t, s, "failed with nil err")
}

func TestErrorf_WithErr_LogsAndExitsOne(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	restore, buf := withCapturedLogs(t)
	defer restore()

	Errorf(errors.New("oops"), "format %s %d", "str", 7)

	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
	s := buf.String()
	assert.Contains(t, s, "\"level\":\"error\"")
	assert.Contains(t, s, "format str 7")
	assert.Contains(t, s, "oops")
}

func TestConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name         string
		function     func()
		expectedCode int
	}{
		{
			name:         "Success exits with code 0",
			function:     func() { Success() },
			expectedCode: 0,
		},
		{
			name:         "SuccessMsg exits with code 0",
			function:     func() { SuccessMsg("test") },
			expectedCode: 0,
		},
		{
			name:         "Successf exits with code 0",
			function:     func() { Successf("test %s", "ok") },
			expectedCode: 0,
		},
		{
			name:         "ErrorMsg exits with code 1",
			function:     func() { ErrorMsg("test") },
			expectedCode: 1,
		},
		{
			name:         "ErrorMsgf exits with code 1",
			function:     func() { ErrorMsgf("test %s", "err") },
			expectedCode: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockHandler := &MockExitHandler{}
			setExitHandler(mockHandler)
			defer resetExitHandler()

			test.function()

			assert.True(t, mockHandler.ExitCalled)
			assert.Equal(t, test.expectedCode, mockHandler.LastExitCode)
			assert.Equal(t, 1, mockHandler.ExitCount)
		})
	}
}
