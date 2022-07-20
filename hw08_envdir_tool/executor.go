package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for name, val := range env {
		if val == "" {
			os.Unsetenv(name)
		} else {
			os.Setenv(name, val)
		}
	}

	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return exec.ExitError{}.ExitCode()
	}

	return 0
}
