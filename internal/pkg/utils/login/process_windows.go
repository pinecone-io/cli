//go:build windows

package login

import (
	"os/exec"
	"syscall"
)

// detachProcess configures cmd to run in its own process group so it is not
// killed when the parent process exits.
func detachProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
