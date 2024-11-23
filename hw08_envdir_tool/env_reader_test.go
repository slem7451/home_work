package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestReadDir(t *testing.T) {
	t.Run("success testdata/env read", func(t *testing.T) {
		rightEnv := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: `"hello"`, NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Equal(t, len(rightEnv), len(env))

		for k, v := range env {
			require.Equal(t, rightEnv[k], v)
		}
	})

	t.Run("success envdir read", func(t *testing.T) {
		err := os.Mkdir("env", 0)
		require.NoError(t, err)

		createEnv := map[string]string{
			"TESTMULTILINE":       "success\nerr\nerror",
			"TESTTAB":             "\t\t\t",
			"TESTSPACE":           "        ",
			"TESTVALSPACEANDTAB1": "value \t",
			"TESTVALSPACEANDTAB2": "value\t ",
			"TESTTERMINALZEROES":  "\x00\x00\x00\x00",
		}

		for k, v := range createEnv {
			file, err := os.Create("./env/" + k)
			require.NoError(t, err)
			file.WriteString(v)
			file.Close()
		}

		rightEnv := Environment{
			"TESTMULTILINE":       EnvValue{Value: "success", NeedRemove: false},
			"TESTTAB":             EnvValue{Value: "", NeedRemove: true},
			"TESTSPACE":           EnvValue{Value: "", NeedRemove: true},
			"TESTVALSPACEANDTAB1": EnvValue{Value: "value", NeedRemove: false},
			"TESTVALSPACEANDTAB2": EnvValue{Value: "value", NeedRemove: false},
			"TESTTERMINALZEROES":  EnvValue{Value: "\n\n\n\n", NeedRemove: false},
		}

		env, err := ReadDir("./env")
		require.NoError(t, err)
		require.Equal(t, len(rightEnv), len(env))

		for k, v := range env {
			require.Equal(t, rightEnv[k], v)
		}

		err = os.RemoveAll("env")
		require.NoError(t, err)
	})

	t.Run("success empty envdir read", func(t *testing.T) {
		err := os.Mkdir("env", 0)
		require.NoError(t, err)

		env, err := ReadDir("./env")
		require.NoError(t, err)
		require.Equal(t, 0, len(env))

		err = os.RemoveAll("env")
		require.NoError(t, err)
	})

	t.Run("error file env name", func(t *testing.T) {
		err := os.Mkdir("env", 0)
		require.NoError(t, err)

		file, err := os.Create("./env/WRO=NG")
		require.NoError(t, err)
		file.Close()

		_, err = ReadDir("./env")
		require.ErrorIs(t, ErrNameOfEnvFile, err)

		err = os.RemoveAll("env")
		require.NoError(t, err)
	})

	t.Run("error directory", func(t *testing.T) {
		_, err := ReadDir("./env")
		require.Error(t, err)
	})
}
