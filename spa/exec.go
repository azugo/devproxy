//go:build !windows

package spa

import (
	"context"
	"os/exec"
	"syscall"
)

func newCommand(ctx context.Context, path string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}

func killCommand(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	// Using syscall as Process.Kill does not kill child processes
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		return err
	}
	return nil
}
