//go:build windows

package spa

import (
	"context"
	"os"
	"os/exec"
)

func newCommand(ctx context.Context, path string, args ...string) *exec.Cmd {
	a := make([]string, 0, len(args)+2)
	a = append(a, "/c", path)
	a = append(a, args...)
	return exec.CommandContext(ctx, "cmd.exe", a...)
}

func killCommand(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	p, err := os.FindProcess(int(cmd.Process.Pid))
	if err != nil {
		return err
	}
	return p.Kill()
}
