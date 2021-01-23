package utils

import (
	"context"
	"os"
	"os/exec"
	"syscall"
)

func Exec(path, name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Dir = path
	return c.Run()
}

func ExecWithPrintStd(ctx context.Context, path, name string, args ...string) (command *exec.Cmd, err error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = path
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd, cmd.Start()
}

func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
