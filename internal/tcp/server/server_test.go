package server

import (
	"context"
	"io"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	const maxConnections = 5
	const maxMessageSize = 100
	server, err := NewServer(
		"127.0.0.1:0",
		maxMessageSize,
		slog.New(slog.NewJSONHandler(io.Discard, nil)),
		WithIdleTimeout(time.Second),
		WithMaxConnections(maxConnections))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	go server.Run(ctx, func(b []byte) ([]byte, error) { // nolint
		return []byte("OK"), nil
	})
	time.Sleep(100 * time.Millisecond)

	address := server.listener.Addr().String()
	command := []byte("set name Daniil")

	t.Run("success", func(t *testing.T) {
		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		t.Cleanup(func() { conn.Close() }) // nolint

		_, err = common.Write(conn, command)
		require.NoError(t, err)

		buffer := make([]byte, 1024)
		n, err := common.Read(conn, buffer)
		require.NoError(t, err)

		assert.Equal(t, "OK", string(buffer[:n]))
	})

	t.Run("exceed client limit", func(t *testing.T) {
		for range maxConnections {
			conn, err := net.Dial("tcp", address)
			require.NoError(t, err)
			t.Cleanup(func() { conn.Close() }) // nolint
		}

		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		t.Cleanup(func() { conn.Close() }) // nolint

		_, err = common.Write(conn, command)
		require.NoError(t, err)

		err = conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		require.NoError(t, err)

		_, err = conn.Read(make([]byte, 10))
		assert.ErrorIs(t, err, os.ErrDeadlineExceeded)
	})

	t.Run("exceed read deadline", func(t *testing.T) {
		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		t.Cleanup(func() { conn.Close() }) // nolint

		time.Sleep(1100 * time.Millisecond)

		_, err = conn.Read(make([]byte, 10))
		assert.ErrorIs(t, err, io.EOF)
	})
}
