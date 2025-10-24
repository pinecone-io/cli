package oauth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTokenErrorKind_String(t *testing.T) {
	tests := []struct {
		kind     TokenErrorKind
		expected string
	}{
		{TokenErrUnknown, "unknown"},
		{TokenErrInvalidRequest, "invalid_request"},
		{TokenErrInvalidClient, "invalid_client"},
		{TokenErrSessionExpired, "session_expired"},
		{TokenErrUnauthorizedClient, "unauthorized_client"},
		{TokenErrUnsupportedGrantType, "unsupported_grant_type"},
		{TokenErrInvalidScope, "invalid_scope"},
		{TokenErrRateLimited, "rate_limited"},
		{TokenErrAuthServerIssue, "auth_server_issue"},
		{TokenErrorKind(999), "unknown"}, // Unknown kind
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.kind.String()
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestClassifyTokenErrorKind(t *testing.T) {
	tests := []struct {
		name     string
		error    *TokenError
		expected TokenErrorKind
	}{
		{
			name: "invalid request",
			error: &TokenError{
				ErrorCode:  ErrInvalidRequest,
				HTTPStatus: 400,
			},
			expected: TokenErrInvalidRequest,
		},
		{
			name: "invalid client",
			error: &TokenError{
				ErrorCode:  ErrInvalidClient,
				HTTPStatus: 401,
			},
			expected: TokenErrInvalidClient,
		},
		{
			name: "invalid grant",
			error: &TokenError{
				ErrorCode:  ErrInvalidGrant,
				HTTPStatus: 403,
			},
			expected: TokenErrSessionExpired,
		},
		{
			name: "unauthorized client",
			error: &TokenError{
				ErrorCode:  ErrUnauthorizedClient,
				HTTPStatus: 401,
			},
			expected: TokenErrUnauthorizedClient,
		},
		{
			name: "unsupported grant type",
			error: &TokenError{
				ErrorCode:  ErrUnsupportedGrantType,
				HTTPStatus: 405,
			},
			expected: TokenErrUnsupportedGrantType,
		},
		{
			name: "invalid scope",
			error: &TokenError{
				ErrorCode:  ErrInvalidScope,
				HTTPStatus: 406,
			},
			expected: TokenErrInvalidScope,
		},
		{
			name: "rate limited",
			error: &TokenError{
				HTTPStatus: 429,
			},
			expected: TokenErrRateLimited,
		},
		{
			name: "auth server issue",
			error: &TokenError{
				HTTPStatus: 500,
			},
			expected: TokenErrAuthServerIssue,
		},
		{
			name: "unknown",
			error: &TokenError{
				HTTPStatus: 404,
			},
			expected: TokenErrUnknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := classifyTokenErrorKind(test.error)
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

// Tests used to validate TokenError.Error() and TokenError.UserMessage()
// Error() basically wraps UserMessage() except in the nil case
var tokenErrTests = []struct {
	name     string
	error    *TokenError
	expected string
}{
	{
		name: "invalid request",
		error: &TokenError{
			Kind: TokenErrInvalidRequest,
		},
		expected: "The authentication server rejected the request. Please try again later.",
	},
	{
		name: "invalid client",
		error: &TokenError{
			Kind: TokenErrInvalidClient,
		},
		expected: "The authentication server rejected the client. Please try again later.",
	},
	{
		name: "session expired",
		error: &TokenError{
			Kind: TokenErrSessionExpired,
		},
		expected: "Your session has expired. Run `pc auth login` to sign in again.",
	},
	{
		name: "unauthorized client",
		error: &TokenError{
			Kind: TokenErrUnauthorizedClient,
		},
		expected: "The authentication server rejected the client. Please try again later.",
	},
	{
		name: "unsupported grant type",
		error: &TokenError{
			Kind: TokenErrUnsupportedGrantType,
		},
		expected: "The authentication server does not support the grant type. Please try again later.",
	},
	{
		name: "invalid scope",
		error: &TokenError{
			Kind: TokenErrInvalidScope,
		},
		expected: "The authentication server rejected the scope. Please try again later.",
	},
	{
		name: "rate limited",
		error: &TokenError{
			Kind: TokenErrRateLimited,
		},
		expected: "Too many requests. Please wait a moment and try again.",
	},
	{
		name: "auth server issue",
		error: &TokenError{
			Kind: TokenErrAuthServerIssue,
		},
		expected: "An error occurred with the authentication server. Please try again later.",
	},
	{
		name: "unknown",
		error: &TokenError{
			Kind: TokenErrUnknown,
		},
		expected: "An unknown error occurred. Please try again later.",
	},
}

func TestTokenError_Error(t *testing.T) {
	baseTests := []struct {
		name     string
		error    *TokenError
		expected string
	}{
		{
			name:     "nil error",
			error:    nil,
			expected: "<nil>",
		},
	}
	tests := append(baseTests, tokenErrTests...)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.error.Error()
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestTokenError_UserMessage(t *testing.T) {
	baseTests := []struct {
		name     string
		error    *TokenError
		expected string
	}{
		{
			name:     "nil error",
			error:    nil,
			expected: "authentication failed",
		},
	}
	tests := append(baseTests, tokenErrTests...)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.error.UserMessage()
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestTokenError_Format(t *testing.T) {
	err := &TokenError{
		Kind:        TokenErrSessionExpired,
		HTTPStatus:  403,
		Description: "Token expired",
		ErrorCode:   ErrInvalidGrant,
		ErrorURI:    "https://example.com/error",
		RawBody:     `{"error":"invalid_grant","error_description":"Unknown or invalid refresh token."}`,
		Operation:   OpRefresh,
	}

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "string format",
			format:   "%s",
			expected: "Your session has expired. Run `pc auth login` to sign in again.",
		},
		{
			name:     "quoted format",
			format:   "%q",
			expected: "\"Your session has expired. Run `pc auth login` to sign in again.\"",
		},
		{
			name:     "verbose format",
			format:   "%+v",
			expected: "token error: kind=session_expired http_status=403 error_code=invalid_grant operation=refresh_token description=\"Token expired\" error_uri=\"https://example.com/error\" raw_body=\"{\\\"error\\\":\\\"invalid_grant\\\",\\\"error_description\\\":\\\"Unknown or invalid refresh token.\\\"}\"",
		},
		{
			name:     "default format",
			format:   "%v",
			expected: "Your session has expired. Run `pc auth login` to sign in again.",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := fmt.Sprintf(test.format, err)
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestNewTokenErrorFromResponse(t *testing.T) {
	tests := []struct {
		name                string
		operation           TokenOperation
		statusCode          int
		responseBody        string
		expectedKind        TokenErrorKind
		expectedStatus      int
		expectedCode        OAuthErrorCode
		expectedDescription string
		expectedErrorURI    string
		expectedRawBody     string
		expectedOperation   TokenOperation
	}{
		{
			name:                "nil response",
			operation:           OpRefresh,
			statusCode:          0,
			responseBody:        "",
			expectedKind:        TokenErrUnknown,
			expectedStatus:      0,
			expectedCode:        "",
			expectedDescription: "",
			expectedErrorURI:    "",
			expectedRawBody:     "",
			expectedOperation:   OpRefresh,
		},
		{
			name:                "invalid grant response",
			operation:           OpRefresh,
			statusCode:          403,
			responseBody:        `{"error":"invalid_grant","error_description":"Unknown or invalid refresh token.","error_uri":"https://example.com/error"}`,
			expectedKind:        TokenErrSessionExpired,
			expectedStatus:      403,
			expectedCode:        ErrInvalidGrant,
			expectedDescription: "Unknown or invalid refresh token.",
			expectedErrorURI:    "https://example.com/error",
			expectedRawBody:     `{"error":"invalid_grant","error_description":"Unknown or invalid refresh token.","error_uri":"https://example.com/error"}`,
			expectedOperation:   OpRefresh,
		},
		{
			name:                "rate limited response",
			operation:           OpRefresh,
			statusCode:          429,
			responseBody:        `{"error":"rate_limited"}`,
			expectedKind:        TokenErrRateLimited,
			expectedStatus:      429,
			expectedCode:        OAuthErrorCode("rate_limited"),
			expectedDescription: "",
			expectedErrorURI:    "",
			expectedRawBody:     `{"error":"rate_limited"}`,
			expectedOperation:   OpRefresh,
		},
		{
			name:                "server error response",
			operation:           OpExchangeCode,
			statusCode:          500,
			responseBody:        `{"error":"server_error","error_description":"Internal server error"}`,
			expectedKind:        TokenErrAuthServerIssue,
			expectedStatus:      500,
			expectedCode:        OAuthErrorCode("server_error"),
			expectedDescription: "Internal server error",
			expectedErrorURI:    "",
			expectedRawBody:     `{"error":"server_error","error_description":"Internal server error"}`,
			expectedOperation:   OpExchangeCode,
		},
		{
			name:                "malformed JSON response",
			operation:           OpRefresh,
			statusCode:          400,
			responseBody:        `{"error":"invalid_resp}`,
			expectedKind:        TokenErrInvalidRequest,
			expectedStatus:      400,
			expectedCode:        "",
			expectedDescription: "",
			expectedErrorURI:    "",
			expectedRawBody:     `{"error":"invalid_resp}`,
			expectedOperation:   OpRefresh,
		},
		{
			name:                "empty body response",
			operation:           OpExchangeCode,
			statusCode:          406,
			responseBody:        "",
			expectedKind:        TokenErrUnknown,
			expectedStatus:      406,
			expectedCode:        OAuthErrorCode(""),
			expectedDescription: "",
			expectedErrorURI:    "",
			expectedRawBody:     "",
			expectedOperation:   OpExchangeCode,
		},
		{
			name:                "truncated raw body",
			operation:           OpRefresh,
			statusCode:          400,
			responseBody:        strings.Repeat("x", 2049),
			expectedKind:        TokenErrInvalidRequest,
			expectedStatus:      400,
			expectedCode:        "",
			expectedDescription: "",
			expectedErrorURI:    "",
			expectedRawBody:     strings.Repeat("x", 2048) + "...(truncated)",
			expectedOperation:   OpRefresh,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var resp *http.Response

			if test.name == "nil response" {
				resp = nil
			} else {
				// Create a mock HTTP response
				resp = &http.Response{
					StatusCode: test.statusCode,
					Body:       http.NoBody, // body handled below in httptest.NewServer
				}
			}

			if resp != nil {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(test.statusCode)
					_, _ = w.Write([]byte(test.responseBody))
				}))
				defer server.Close()

				client := &http.Client{}
				req, _ := http.NewRequest("GET", server.URL, nil)
				resp, _ = client.Do(req)
			}

			got := NewTokenErrorFromResponse(test.operation, resp)
			if got.Kind != test.expectedKind {
				t.Errorf("expected kind %q, got %q", test.expectedKind, got.Kind)
			}
			if got.HTTPStatus != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, got.HTTPStatus)
			}
			if got.ErrorCode != test.expectedCode {
				t.Errorf("expected code %q, got %q", test.expectedCode, got.ErrorCode)
			}
			if got.Description != test.expectedDescription {
				t.Errorf("expected description %q, got %q", test.expectedDescription, got.Description)
			}
			if got.ErrorURI != test.expectedErrorURI {
				t.Errorf("expected error URI %q, got %q", test.expectedErrorURI, got.ErrorURI)
			}
			if got.RawBody != test.expectedRawBody {
				t.Errorf("expected raw body %q, got %q", test.expectedRawBody, got.RawBody)
			}
			if got.Operation != test.expectedOperation {
				t.Errorf("expected operation %q, got %q", test.expectedOperation, got.Operation)
			}
		})
	}
}
