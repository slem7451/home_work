package main

import (
	"errors"
	"os"
	"os/exec"
	"slices"
	"strings"
	"fmt"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)

	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	command.Env = os.Environ()

	for k, v := range env {
		if v.NeedRemove {
			idx := slices.IndexFunc(command.Env, func(item string) bool {
				return strings.HasPrefix(item, k + "=")
			})

			if idx > 0 {
				command.Env = append(command.Env[:idx], command.Env[idx+1:]...)
			}
		} else {
			command.Env = append(command.Env, k + "=" + v.Value)
		}
	}

	err := command.Run()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return exit.ExitCode()
		}
		fmt.Fprint(os.Stderr, err)

		return 1
	}

	return 0
}
