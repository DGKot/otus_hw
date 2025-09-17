package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	com := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout
	com.Stderr = os.Stderr
	com.Env = append(os.Environ(), env.ToSliceString()...)
	err := com.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		fmt.Fprintln(os.Stderr, "error with run command:", err)
		return 1
	}

	return 0
}
