package goup

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func Restart() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}

	// Windows does not support exec syscall.
	if runtime.GOOS == "windows" {
		cmd := exec.Command(self, os.Args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()

		if err := cmd.Start(); err != nil {
			return err
		}
	} else {
		if err := syscall.Exec(self, os.Args, os.Environ()); err != nil {
			return err
		}
	}

	os.Exit(0)
	return nil
}
