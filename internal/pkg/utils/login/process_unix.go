//go:build !windows

package login

import (
	"os/exec"
	"syscall"
)

// detachProcess configures cmd to run as a new session leader so it is not
// killed when the parent process exits or its controlling terminal closes.
func detachProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
