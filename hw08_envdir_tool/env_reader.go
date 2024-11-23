package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrNameOfEnvFile = errors.New("invalid env file name")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	readedDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, v := range readedDir {
		if v.IsDir() {
			continue
		}

		if strings.Contains(v.Name(), "=") {
			return nil, ErrNameOfEnvFile
		}

		file, err := os.Open(filepath.Join(dir, v.Name()))
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(file)
		envVal, _, err := reader.ReadLine()
		file.Close()

		if err != io.EOF && err != nil {
			return nil, err
		}

		envString := string(envVal)

		envString = strings.ReplaceAll(envString, "\x00", "\n")
		envString = strings.TrimRight(envString, "\t ")

		env[v.Name()] = EnvValue{Value: envString, NeedRemove: len(envString) == 0}
	}

	return env, nil
}
