package login

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
)

type SessionState struct {
	SessionId string    `json:"session_id"`
	CSRFState string    `json:"csrf_state"`
	AuthURL   string    `json:"auth_url"`
	OrgId     *string   `json:"org_id,omitempty"`
	// SSOConnection is set on the second-round SSO session. A non-nil value
	// means this session was started specifically for SSO enforcement, so the
	// completion handler should skip the SSO check and emit "authenticated".
	SSOConnection *string   `json:"sso_connection,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type SessionResult struct {
	SessionId   string    `json:"session_id"`
	Status      string    `json:"status"` // "success" or "error"
	Error       string    `json:"error,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
}

const sessionMaxAge = 5 * time.Minute

func sessionsDir() (string, error) {
	dir := filepath.Join(configuration.ConfigDirPath(), "sessions")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("error creating sessions directory: %w", err)
	}
	return dir, nil
}

func sessionStatePath(sessionId string) (string, error) {
	dir, err := sessionsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, sessionId+".json"), nil
}

func sessionResultPath(sessionId string) (string, error) {
	dir, err := sessionsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, sessionId+".result"), nil
}

func writeSessionState(s SessionState) error {
	path, err := sessionStatePath(s.SessionId)
	if err != nil {
		return err
	}
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("error marshaling session state: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

func ReadSessionState(sessionId string) (*SessionState, error) {
	path, err := sessionStatePath(sessionId)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading session state: %w", err)
	}
	var s SessionState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error unmarshaling session state: %w", err)
	}
	return &s, nil
}

func WriteSessionResult(r SessionResult) error {
	path, err := sessionResultPath(r.SessionId)
	if err != nil {
		return err
	}
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshaling session result: %w", err)
	}
	// Write to a temp file in the same directory, then rename for atomicity.
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".result-*.tmp")
	if err != nil {
		return fmt.Errorf("error creating temp result file: %w", err)
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("error writing temp result file: %w", err)
	}
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("error setting permissions on temp result file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("error closing temp result file: %w", err)
	}
	return os.Rename(tmpPath, path)
}

func readSessionResult(sessionId string) (*SessionResult, error) {
	path, err := sessionResultPath(sessionId)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // not complete yet
		}
		return nil, fmt.Errorf("error reading session result: %w", err)
	}
	var r SessionResult
	if err := json.Unmarshal(data, &r); err != nil {
		// Partial write in progress — treat as not ready yet.
		return nil, nil
	}
	return &r, nil
}

func CleanupSession(sessionId string) {
	statePath, _ := sessionStatePath(sessionId)
	resultPath, _ := sessionResultPath(sessionId)
	_ = os.Remove(statePath)
	_ = os.Remove(resultPath)
}

// findResumableSession looks for a recent pending session that can be resumed.
// Returns the session state and its result (nil if still pending), or nil if
// no resumable session exists.
func findResumableSession() (*SessionState, *SessionResult, error) {
	dir, err := sessionsDir()
	if err != nil {
		return nil, nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading sessions directory: %w", err)
	}

	var newest *SessionState
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		id := e.Name()[:len(e.Name())-len(".json")]
		s, err := ReadSessionState(id)
		if err != nil {
			continue
		}
		if time.Since(s.CreatedAt) > sessionMaxAge {
			// Stale — clean it up
			CleanupSession(s.SessionId)
			continue
		}
		if newest == nil || s.CreatedAt.After(newest.CreatedAt) {
			newest = s
		}
	}

	if newest == nil {
		return nil, nil, nil
	}

	result, err := readSessionResult(newest.SessionId)
	if err != nil {
		return nil, nil, err
	}

	return newest, result, nil
}

func newSessionId() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
