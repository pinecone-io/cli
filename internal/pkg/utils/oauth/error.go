package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OAuthErrorCode string

// OAuth2 Error codes defined in RFC6749: https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
// These codes denote what kind of error the oauth server returned, and we use them to try and derive TokenErrorKind
// along with HTTP status code in the response.
const (
	ErrInvalidRequest       OAuthErrorCode = "invalid_request"
	ErrInvalidClient        OAuthErrorCode = "invalid_client"
	ErrInvalidGrant         OAuthErrorCode = "invalid_grant"
	ErrUnauthorizedClient   OAuthErrorCode = "unauthorized_client"
	ErrUnsupportedGrantType OAuthErrorCode = "unsupported_grant_type"
	ErrInvalidScope         OAuthErrorCode = "invalid_scope"
)

type TokenErrorKind int

const (
	TokenErrUnknown TokenErrorKind = iota
	TokenErrInvalidRequest
	TokenErrInvalidClient
	TokenErrSessionExpired
	TokenErrInvalidGrant
	TokenErrUnauthorizedClient
	TokenErrUnsupportedGrantType
	TokenErrInvalidScope
	TokenErrRateLimited
	TokenErrAuthServerIssue
)

func (k TokenErrorKind) String() string {
	switch k {
	case TokenErrInvalidRequest:
		return "invalid_request"
	case TokenErrInvalidClient:
		return "invalid_client"
	case TokenErrSessionExpired:
		return "session_expired"
	case TokenErrInvalidGrant:
		return "invalid_grant"
	case TokenErrUnauthorizedClient:
		return "unauthorized_client"
	case TokenErrUnsupportedGrantType:
		return "unsupported_grant_type"
	case TokenErrInvalidScope:
		return "invalid_scope"
	case TokenErrRateLimited:
		return "rate_limited"
	case TokenErrAuthServerIssue:
		return "auth_server_issue"
	default:
		return "unknown"
	}
}

// Tagging the operation the error originated from
type TokenOperation string

const (
	OpRefresh      TokenOperation = "refresh_token"
	OpExchangeCode TokenOperation = "exchange_code"
)

type TokenError struct {
	Kind        TokenErrorKind
	HTTPStatus  int
	ErrorCode   OAuthErrorCode
	Description string
	ErrorURI    string
	RawBody     string
	Operation   TokenOperation
}

// Returns the user-facing error message for TokenError
// Format("%+v") is used for full error details
func (e *TokenError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.UserMessage()
}

func (e *TokenError) UserMessage() string {
	if e == nil {
		return "authentication failed"
	}
	switch e.Kind {
	case TokenErrSessionExpired:
		return "Your session has expired. Run `pc auth login` to sign in again."
	case TokenErrRateLimited:
		return "Too many requests. Please wait a moment and try again."
	case TokenErrAuthServerIssue:
		return "An error occurred with the authentication server. Please try again later."
	case TokenErrInvalidClient:
		return "The authentication server rejected the client. Please try again later."
	case TokenErrInvalidRequest:
		return "The authentication server rejected the request. Please try again later."
	case TokenErrUnknown:
		return "An unknown error occurred. Please try again later."
	default:
		if e.Description != "" {
			return "Authentication failed: " + e.Description
		}
		return "Authentication failed."
	}
}

// Handles custom formatting of TokenError for use with fmt.Printf, fmt.Println, etc
func (e *TokenError) Format(s fmt.State, verb rune) {
	if e == nil {
		_, _ = io.WriteString(s, "<nil>")
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			// full details
			fmt.Fprintf(s, "token error: kind=%s http_status=%d error_code=%s operation=%s",
				e.Kind, e.HTTPStatus, e.ErrorCode, e.Operation)
			if e.Description != "" {
				fmt.Fprintf(s, " description=%q", e.Description)
			}
			if e.ErrorURI != "" {
				fmt.Fprintf(s, " error_uri=%q", e.ErrorURI)
			}
			if e.RawBody != "" {
				fmt.Fprintf(s, " raw_body=%q", e.RawBody)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.UserMessage())
	case 'q':
		fmt.Fprintf(s, "%q", e.UserMessage())
	default:
		_, _ = io.WriteString(s, e.UserMessage())
	}
}

type oauthErrPayload struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}

// Parses a non-2xx OAuth token endpoint response body into a TokenError
func NewTokenErrorFromResponse(op TokenOperation, resp *http.Response) *TokenError {
	if resp == nil {
		return &TokenError{
			Kind:      TokenErrUnknown,
			Operation: op,
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB cap
	raw := string(body)
	if len(raw) > 2048 {
		raw = raw[:2048] + "...(truncated)"
	}

	var p oauthErrPayload
	_ = json.Unmarshal(body, &p)

	e := &TokenError{
		HTTPStatus:  resp.StatusCode,
		ErrorCode:   OAuthErrorCode(p.Error),
		Description: p.ErrorDescription,
		ErrorURI:    p.ErrorURI,
		RawBody:     raw,
		Operation:   op,
	}
	e.Kind = classifyTokenErrorKind(e)
	return e
}

func classifyTokenErrorKind(err *TokenError) TokenErrorKind {
	switch err.ErrorCode {
	case ErrInvalidGrant:
		return TokenErrSessionExpired
	case ErrInvalidClient:
		return TokenErrInvalidClient
	case ErrInvalidRequest:
		return TokenErrInvalidRequest
	case ErrUnauthorizedClient:
		return TokenErrUnauthorizedClient
	case ErrUnsupportedGrantType:
		return TokenErrUnsupportedGrantType
	case ErrInvalidScope:
		return TokenErrInvalidScope
	default:
		// Unknown oauth error code, fall back to HTTP semantics
		if err.HTTPStatus == 429 {
			return TokenErrRateLimited
		} else if err.HTTPStatus >= 500 {
			return TokenErrAuthServerIssue
		} else if err.HTTPStatus == 400 {
			return TokenErrInvalidRequest
		} else {
			return TokenErrUnknown
		}
	}
}
