package exit

import (
	"testing"

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

func TestSuccessExitEventMethods(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	successEvent := Success()

	// Msg
	successEvent.Msg("test")
	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 0, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)

	mockHandler.Reset()

	// Msgf
	successEvent.Msgf("test %s", "foo")
	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 0, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)

	mockHandler.Reset()

	// Send
	successEvent.Send()
	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 0, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
}

func TestErrorExitEventMethods(t *testing.T) {
	mockHandler := &MockExitHandler{}
	setExitHandler(mockHandler)
	defer resetExitHandler()

	errorEvent := Error()

	// Msg
	errorEvent.Msg("test")
	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)

	mockHandler.Reset()

	// Msgf
	errorEvent.Msg("test")
	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)

	mockHandler.Reset()

	// Msgf
	errorEvent.Msgf("test %s", "foo")
	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)

	mockHandler.Reset()

	// Send
	errorEvent.Send()
	assert.True(t, mockHandler.ExitCalled)
	assert.Equal(t, 1, mockHandler.LastExitCode)
	assert.Equal(t, 1, mockHandler.ExitCount)
}

func TestConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name         string
		function     func()
		expectedCode int
	}{
		{
			name:         "SuccessMsg exits with code 0",
			function:     func() { SuccessMsg("test") },
			expectedCode: 0,
		},
		{
			name:         "ErrorMsg exits with code 1",
			function:     func() { ErrorMsg("test") },
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
