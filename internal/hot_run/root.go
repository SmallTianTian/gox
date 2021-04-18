package hot_run

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/SmallTianTian/go-tools/slice"
	"github.com/fsnotify/fsnotify"
)

const delayTime = 5 * time.Second

var command *exec.Cmd
var ignoreFile = []string{"bin", "vendor", "go.mod", "go.sum"}

func ClearRun() {
	closeCommand(command)
}

func HotRun(commands []string) {
	comm, args := polishingCommand(commands)
	var err error
	c := fileChange()
	for {
		command, err = buildCommand(comm, args)
		select {
		case <-syncWait(command):
			utils.MustTrue(command.ProcessState.ExitCode() <= 0, "Exit with error.")
			utils.MustNotError(err)
		case <-c:
			closeCommand(command)
		}
	}
}

func syncWait(command *exec.Cmd) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		command.Wait()
		c <- struct{}{}
	}()
	return c
}

// https://stackoverflow.com/a/29552044
func closeCommand(command *exec.Cmd) {
	if command != nil {
		pgid, err := syscall.Getpgid(command.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, 15) // note the minus sign
		}
	}
}

func polishingCommand(comms []string) (comm string, args []string) {
	if len(comms) == 0 {
		return "make", []string{"run"}
	}

	if !utils.CheckCommandExists(comms[0]) {
		return "make", comms
	}
	return comms[0], comms[1:]
}

func buildCommand(comm string, args []string) (command *exec.Cmd, err error) {
	path, _ := os.Getwd()
	if command, err = utils.ExecWithPrintStd(context.Background(), path, comm, args...); err != nil {
		return
	}
	return
}

func canReboot() bool {
	// 1. rebuild go mod
	path, err := os.Getwd()
	utils.MustNotError(err)
	if flag := utils.GoModRebuild(path); !flag {
		return false
	}

	// 2. rebuild exe
	if err = utils.Exec(path, "make", "build"); err != nil {
		return false
	}

	// 3. check config(will do)
	return true
}

func fileChange() <-chan struct{} {
	c := make(chan struct{})

	watcher, err := fsnotify.NewWatcher()
	utils.MustNotError(err)
	addFile2Watcher(watcher)

	go func() {
		defer close(c)
		defer watcher.Close()
		for {
			time.Sleep(delayTime)
			cleanChan(watcher.Events)
			select {
			case event := <-watcher.Events:
				fmt.Println(event)
				if canReboot() {
					c <- struct{}{}
				}
			case err := <-watcher.Errors:
				utils.MustNotError(err)
			}
		}
	}()
	return c
}

func cleanChan(c <-chan fsnotify.Event) {
	for {
		select {
		case <-c:
			fmt.Println("a ...interface{}")
		default:
			return
		}
	}
}

func addFile2Watcher(watcher *fsnotify.Watcher) {
	utils.MustNotError(filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		paths := strings.Split(path, string(os.PathSeparator))
		first := paths[0]
		last := paths[len(paths)-1]

		// ignore hidden file
		if first[0] == '.' || last[0] == '.' {
			return nil
		}

		if slice.StringInSlice(ignoreFile, first) || slice.StringInSlice(ignoreFile, last) {
			return nil
		}
		return watcher.Add(path)
	}))
}
