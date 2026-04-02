//go:build windows

package login

import (
	"os/exec"
	"syscall"
)

// detachProcess configures cmd to run detached from the parent's console session.
// CREATE_NEW_PROCESS_GROUP prevents Ctrl+C propagation; DETACHED_PROCESS (0x00000008)
// detaches the child from the parent's console so it survives terminal close.
func detachProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008,
	}
}
