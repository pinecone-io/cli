//go:build e2e

package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Generic waiter for conditions
func WaitUntil(max time.Duration, initialBackoff time.Duration, cond func() (bool, error)) error {
	deadline := time.Now().Add(max)
	backoff := initialBackoff
	for {
		ok, err := cond()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		if time.Now().After(deadline) {
			return errors.New("timeout")
		}
		time.Sleep(backoff)
		if backoff < 15*time.Second {
			backoff += 2 * time.Second
		}
	}
}

type JSONRunner interface {
	RunJSON(out any, args ...string) (string, error)
}

func WaitForIndexReady(c JSONRunner, name string, max time.Duration) error {
	return WaitUntil(max, 5*time.Second, func() (bool, error) {
		var idx struct {
			Status struct {
				State string `json:"state"`
			} `json:"status"`
		}
		stdout, err := c.RunJSON(&idx, "index", "describe", "--name", name)
		if err != nil {
			// If index not found yet, keep polling
			if strings.Contains(stdout, "does not exist") {
				return false, nil
			}
			return false, fmt.Errorf("describe failed: %w", err)
		}
		// handle exact or lowercased states
		if strings.EqualFold(idx.Status.State, "ready") {
			return true, nil
		}
		// tolerate unexpected shapes by scanning
		var m map[string]any
		if err := json.Unmarshal([]byte(stdout), &m); err == nil {
			if s, ok := deepLookupString(m, "status", "state"); ok && strings.EqualFold(s, "ready") {
				return true, nil
			}
		}
		return false, nil
	})
}

func deepLookupString(m map[string]any, keys ...string) (string, bool) {
	var cur any = m
	for _, k := range keys {
		asMap, ok := cur.(map[string]any)
		if !ok {
			return "", false
		}
		cur, ok = asMap[k]
		if !ok {
			return "", false
		}
	}
	s, ok := cur.(string)
	return s, ok
}
