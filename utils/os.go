package utils

import (
	"bufio"
	"context"
	"io"
	"os/exec"
)

func Exec(path, name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Dir = path
	return c.Run()
}

func ExecWithOut(ctx context.Context, path, name string, args ...string) (<-chan string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = path
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c := make(chan string, 0)
	go func() {
		defer close(c)
		defer stdout.Close()
		br := bufio.NewReader(stdout)
		for {
			if ctx.Err() != nil {
				return
			}
			lbs, _, e := br.ReadLine()
			if e == nil {
				c <- string(lbs)
				continue
			}
			if e != io.EOF {
				c <- string("Error: " + e.Error())
			}
			return
		}
	}()
	return c, cmd.Run()
}

func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
