package process

import (
	"io"
	"os"
)

type Config struct {
	Binary  string
	Args    []string
	Workdir string
}

func RunProcess(cfg Config) (done, wait func(), err error) {
	command := generateCommand(cfg)
	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	go func() {
		_, _ = io.Copy(os.Stdout, stdoutPipe)
		_ = stdoutPipe.Close()
	}()
	go func() {
		_, _ = io.Copy(os.Stderr, stderrPipe)
		_ = stderrPipe.Close()
	}()
	if err := command.Start(); err != nil {
		return nil, nil, err
	}
	done = func() {
		_ = command.Process.Kill()
	}
	wait = func() {
		_ = command.Wait()
	}
	return done, wait, nil
}
