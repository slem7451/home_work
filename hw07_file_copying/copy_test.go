package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestCopy(t *testing.T) {
	from := "testdata/input.txt"
	to := "out.txt"

	t.Run("limit out of size", func(t *testing.T) {
		err := Copy(from, to, 0, 10000000000)
		require.NoError(t, err)

		in, err := os.ReadFile(from)
		require.NoError(t, err)

		out, err := os.ReadFile(to)
		require.NoError(t, err)

		require.True(t, bytes.Equal(in, out))
	})

	t.Run("full copy", func(t *testing.T) {
		err := Copy(from, to, 0, 0)
		require.NoError(t, err)

		in, err := os.ReadFile(from)
		require.NoError(t, err)

		out, err := os.ReadFile(to)
		require.NoError(t, err)

		require.True(t, bytes.Equal(in, out))
	})

	t.Run("10 limit and 10 offset", func(t *testing.T) {
		err := Copy(from, to, 10, 10)
		require.NoError(t, err)

		out, err := os.ReadFile(to)
		require.NoError(t, err)

		require.True(t, bytes.Equal([]byte("ts\nPackage"), out))
	})

	os.Remove(to)

	t.Run("error 'offset exceeds file size' case", func(t *testing.T) {
		err := Copy(from, to, 100000000, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("from and to the same", func(t *testing.T) {
		path := "test.txt"

		file, err := os.Create(path)
		require.NoError(t, err)

		file.Write([]byte("testtesttesttesttesttest"))
		file.Close()

		err = Copy(path, path, 4, 4)
		require.NoError(t, err)

		out, err := os.ReadFile(path)
		require.NoError(t, err)

		require.True(t, bytes.Equal([]byte("test"), out))

		os.Remove(path)
	})
}
