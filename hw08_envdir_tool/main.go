package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrNeedMoreArgs = errors.New("arguments less then 3")

func main() {
	envDir, cmd := parseArgs()

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	code := RunCmd(cmd, env)
	os.Exit(code)
}

func parseArgs() (string, []string) {
	args := os.Args

	if len(args) < 3 {
		fmt.Fprint(os.Stderr, ErrNeedMoreArgs)
		os.Exit(1)
	}

	return args[1], args[2:]
}
