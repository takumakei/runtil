package app

import (
	"os"
	"os/exec"
	"syscall"
)

type Process struct {
	C   <-chan error
	Pid int
	cmd *exec.Cmd
}

func Start(name string, arg ...string) *Process {
	c := make(chan error)

	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pid := make(chan int)

	go func() {
		if err := cmd.Start(); err != nil {
			c <- err
			return
		}
		pid <- cmd.Process.Pid
		c <- cmd.Wait()
	}()

	return &Process{C: c, Pid: <-pid, cmd: cmd}
}

func (pcs *Process) Kill(sig os.Signal) error {
	return syscall.Kill(pcs.cmd.Process.Pid, sig.(syscall.Signal))
}
