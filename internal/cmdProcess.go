package internal

import (
	"os/exec"
)

type cmdProcess struct {
	cmd *exec.Cmd
}

func (proc *cmdProcess) Start() error {
	return proc.cmd.Start()
}

func (proc *cmdProcess) Kill() error {
	if proc.cmd.Process == nil {
		return nil
	}
	return proc.cmd.Process.Kill()
}

func (proc *cmdProcess) Wait() error {
	if proc.cmd.Process == nil {
		return nil
	}
	return proc.cmd.Wait()
}
