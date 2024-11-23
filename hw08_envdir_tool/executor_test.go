package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestRunCmd(t *testing.T) {
	t.Run("success exec with empty env", func(t *testing.T) {
		code := RunCmd([]string{"ls"}, make(Environment))
		require.Equal(t, 0, code)
	})

	t.Run("success exec with env", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		code := RunCmd([]string{"ls"}, env)
		require.Equal(t, 0, code)
	})
	
	t.Run("error exec with empty env", func(t *testing.T) {
		code := RunCmd([]string{"cd"}, make(Environment))
		require.Equal(t, 1, code)
	})

	t.Run("error exec with env", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		code := RunCmd([]string{"cd"}, env)
		require.Equal(t, 1, code)
	})
}
