package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("server closed", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		defer l.Close()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(l.Addr().String(), time.Second, io.NopCloser(in), out)
		require.NoError(t, client.Connect())

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			con, err := l.Accept()
			require.NoError(t, err)
			con.Close()
		}()

		wg.Wait()

		in.WriteString("1")
		client.Send()
		in.WriteString("1")

		require.NotNil(t, client.Send())
	})

	t.Run("server not found", func(t *testing.T) {
		client := NewTelnetClient("127.0.0.1:0", time.Second, os.Stdin, os.Stdout)
		require.NotNil(t, client.Connect())
	})

	t.Run("client without connect", func(t *testing.T) {
		client := NewTelnetClient("127.0.0.1:0", time.Second, os.Stdin, os.Stdout)
		require.NoError(t, client.Close())
	})

	t.Run("client Ctrl+D", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		defer l.Close()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(l.Addr().String(), time.Second, io.NopCloser(in), out)
		require.NoError(t, client.Connect())

		go func() {
			con, err := l.Accept()
			require.NoError(t, err)

			_, err = io.Copy(io.Discard, con)
			require.NoError(t, err)
		}()

		in.WriteString("1")
		in.WriteString("1")

		in.Write([]byte{})

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			require.NoError(t, client.Send())
		}()

		wg.Wait()

		require.NoError(t, client.Close())
	})
}
