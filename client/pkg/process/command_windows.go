//go:build windows

package process

import (
	"os/exec"
	"path/filepath"
	"syscall"
)

func generateCommand(cfg Config) *exec.Cmd {
	command := exec.Command(cfg.Binary, cfg.Args...)
	command.Dir = filepath.Dir(cfg.Workdir)
	command.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	return command
}
